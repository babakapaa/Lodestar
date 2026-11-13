import 'dart:async';
import 'dart:convert';

import 'package:drift/drift.dart';

import '../../core/offline/client_id.dart';
import 'chat.dart';
import 'chat_entry.dart';
import 'chat_file.dart';
import 'chat_message.dart';
import 'chats_database.dart';
import 'outgoing.dart';

class ChatsStore {
  ChatsStore(this._db);

  static const _meKey = 'me';

  final ChatsDatabase _db;
  final _local = StreamController<void>.broadcast();

  Stream<void> get localChanges => _local.stream;

  Stream<List<Chat>> watchChats() =>
      (_db.select(_db.chatRows)
            ..orderBy([(row) => OrderingTerm.desc(row.changedAt)]))
          .watch()
          .map((rows) => [for (final row in rows) _chat(row)]);

  Stream<Chat?> watchChat(String id) =>
      (_db.select(_db.chatRows)..where((row) => row.id.equals(id)))
          .watchSingleOrNull()
          .map((row) => row == null ? null : _chat(row));

  Stream<bool> watchOlder(String id) =>
      (_db.select(_db.chatRows)..where((row) => row.id.equals(id)))
          .watchSingleOrNull()
          .map((row) => row?.older ?? false);

  Stream<List<ChatEntry>> watchEntries(String chatId) {
    final sent =
        (_db.select(_db.messageRows)
              ..where((row) => row.chatId.equals(chatId))
              ..orderBy([(row) => OrderingTerm.asc(row.seq)]))
            .watch();
    final pending =
        (_db.select(_db.outboxRows)
              ..where((row) => row.chatId.equals(chatId))
              ..orderBy([(row) => OrderingTerm.asc(row.createdAt)]))
            .watch();
    late StreamController<List<ChatEntry>> out;
    List<MessageRow>? messages;
    List<OutboxRow>? outbox;
    StreamSubscription<List<MessageRow>>? a;
    StreamSubscription<List<OutboxRow>>? b;
    void emit() {
      if (messages == null || outbox == null) {
        return;
      }
      final known = <String>{};
      final entries = <ChatEntry>[];
      for (final row in messages!) {
        final message = _message(row.data);
        if (message.clientId != null) {
          known.add(message.clientId!);
        }
        entries.add(ChatEntry.sent(message));
      }
      for (final row in outbox!) {
        if (!known.contains(row.clientId)) {
          entries.add(ChatEntry.pending(_outgoing(row)));
        }
      }
      out.add(entries);
    }

    out = StreamController<List<ChatEntry>>(
      onListen: () {
        a = sent.listen((rows) {
          messages = rows;
          emit();
        }, onError: out.addError);
        b = pending.listen((rows) {
          outbox = rows;
          emit();
        }, onError: out.addError);
      },
      onCancel: () async {
        await a?.cancel();
        await b?.cancel();
      },
    );
    return out.stream;
  }

  Stream<String?> watchMe() =>
      (_db.select(_db.chatValues)..where((row) => row.key.equals(_meKey)))
          .watchSingleOrNull()
          .map((row) => row?.value);

  Future<String?> me() async => (await (_db.select(
    _db.chatValues,
  )..where((row) => row.key.equals(_meKey))).getSingleOrNull())?.value;

  Future<void> setMe(String id) => _db
      .into(_db.chatValues)
      .insertOnConflictUpdate(
        ChatValuesCompanion.insert(key: _meKey, value: id),
      );

  Future<List<Chat>> chats() async => [
    for (final row in await _db.select(_db.chatRows).get()) _chat(row),
  ];

  Future<int> syncedRev(String id) async =>
      (await (_db.select(
        _db.chatRows,
      )..where((row) => row.id.equals(id))).getSingleOrNull())?.syncedRev ??
      -1;

  Future<void> saveChats(List<Chat> chats) => _db.transaction(() async {
    final ids = {for (final chat in chats) chat.id};
    for (final chat in chats) {
      await _saveChat(chat);
    }
    final gone = (await _db.select(_db.chatRows).get())
        .where((row) => !ids.contains(row.id))
        .map((row) => row.id);
    for (final id in gone.toList()) {
      await forget(id);
    }
  });

  Future<void> saveChat(Chat chat) => _saveChat(chat);

  Future<void> forget(String id) => _db.transaction(() async {
    await (_db.delete(
      _db.messageRows,
    )..where((row) => row.chatId.equals(id))).go();
    await (_db.delete(
      _db.outboxRows,
    )..where((row) => row.chatId.equals(id))).go();
    await (_db.delete(_db.chatRows)..where((row) => row.id.equals(id))).go();
  });

  Future<void> saveMessages(
    String chatId,
    List<ChatMessage> messages, {
    int? syncedRev,
    bool? older,
  }) => _db.transaction(() async {
    for (final message in messages) {
      await _db
          .into(_db.messageRows)
          .insertOnConflictUpdate(
            MessageRowsCompanion.insert(
              id: message.id,
              chatId: chatId,
              seq: message.seq,
              data: jsonEncode(message.toJson()),
            ),
          );
      if (message.clientId != null) {
        await (_db.delete(
          _db.outboxRows,
        )..where((row) => row.clientId.equals(message.clientId!))).go();
      }
    }
    if (syncedRev != null || older != null) {
      await (_db.update(
        _db.chatRows,
      )..where((row) => row.id.equals(chatId))).write(
        ChatRowsCompanion(
          syncedRev: syncedRev == null
              ? const Value.absent()
              : Value(syncedRev),
          older: older == null ? const Value.absent() : Value(older),
        ),
      );
    }
  });

  Future<int?> firstSeq(String chatId) async {
    final row =
        await (_db.select(_db.messageRows)
              ..where((row) => row.chatId.equals(chatId))
              ..orderBy([(row) => OrderingTerm.asc(row.seq)])
              ..limit(1))
            .getSingleOrNull();
    return row?.seq;
  }

  Future<void> markRead(String chatId, int seq) async {
    final row = await (_db.select(
      _db.chatRows,
    )..where((row) => row.id.equals(chatId))).getSingleOrNull();
    if (row == null) {
      return;
    }
    await _saveChat(_chat(row).read(seq));
  }

  Future<String> queue(
    String chatId,
    String text, {
    int? replyTo,
    List<String> mentions = const [],
    List<ChatFile> files = const [],
  }) async {
    final id = ClientId.generate();
    await _db
        .into(_db.outboxRows)
        .insert(
          OutboxRowsCompanion.insert(
            clientId: id,
            chatId: chatId,
            body: text,
            replyTo: Value(replyTo),
            mentions: Value(jsonEncode(mentions)),
            files: Value(jsonEncode([for (final file in files) file.toJson()])),
            createdAt: DateTime.now().toUtc(),
          ),
        );
    _local.add(null);
    return id;
  }

  Future<List<Outgoing>> outbox() async => [
    for (final row
        in await (_db.select(_db.outboxRows)
              ..where((row) => row.error.isNull())
              ..orderBy([(row) => OrderingTerm.asc(row.createdAt)]))
            .get())
      _outgoing(row),
  ];

  Future<void> fail(String clientId, String error) =>
      (_db.update(_db.outboxRows)
            ..where((row) => row.clientId.equals(clientId)))
          .write(OutboxRowsCompanion(error: Value(error)));

  Future<void> retry(String clientId) async {
    await (_db.update(_db.outboxRows)
          ..where((row) => row.clientId.equals(clientId)))
        .write(const OutboxRowsCompanion(error: Value(null)));
    _local.add(null);
  }

  Future<void> discard(String clientId) => (_db.delete(
    _db.outboxRows,
  )..where((row) => row.clientId.equals(clientId))).go();

  Future<void> clear() => _db.clear();

  Future<void> dispose() => _local.close();

  Future<void> _saveChat(Chat chat) => _db
      .into(_db.chatRows)
      .insert(
        ChatRowsCompanion.insert(
          id: chat.id,
          data: jsonEncode(chat.toJson()),
          changedAt: chat.changedAt,
        ),
        onConflict: DoUpdate(
          (_) => ChatRowsCompanion(
            data: Value(jsonEncode(chat.toJson())),
            changedAt: Value(chat.changedAt),
          ),
        ),
      );

  static Chat _chat(ChatRow row) =>
      Chat.fromJson(jsonDecode(row.data) as Map<String, Object?>);

  static ChatMessage _message(String data) =>
      ChatMessage.fromJson(jsonDecode(data) as Map<String, Object?>);

  static Outgoing _outgoing(OutboxRow row) => Outgoing(
    clientId: row.clientId,
    chatId: row.chatId,
    text: row.body,
    replyTo: row.replyTo,
    mentions: [
      for (final id in jsonDecode(row.mentions) as List<Object?>) id! as String,
    ],
    files: [
      for (final item in jsonDecode(row.files) as List<Object?>)
        ChatFile.fromJson(item! as Map<String, Object?>),
    ],
    createdAt: row.createdAt,
    error: row.error,
  );
}
