import '../../core/api/api_exception.dart';
import 'chat.dart';
import 'chats_api.dart';
import 'chats_store.dart';

class ChatsSync {
  const ChatsSync(this._api, this._store);

  final ChatsApi _api;
  final ChatsStore _store;

  static bool offline(ApiException error) =>
      error.code == ApiException.unreachable ||
      error.statusCode == 401 ||
      error.statusCode >= 500;

  Future<void> run() async {
    if (await _store.me() == null) {
      await _store.setMe(await _api.me());
    }
    await flush();
    final remote = await _api.list();
    await _store.saveChats(remote);
    for (final chat in remote) {
      await pull(chat);
    }
  }

  Future<void> flush() async {
    for (final message in await _store.outbox()) {
      try {
        final sent = await _api.send(message);
        await _store.saveMessages(message.chatId, [sent]);
      } on ApiException catch (error) {
        if (offline(error)) {
          rethrow;
        }
        await _store.fail(message.clientId, error.code);
      }
    }
  }

  Future<void> pull(Chat chat) async {
    var synced = await _store.syncedRev(chat.id);
    if (synced < 0) {
      final page = await _api.latest(chat.id);
      await _store.saveMessages(
        chat.id,
        page.messages,
        syncedRev: chat.rev,
        older: page.more,
      );
      return;
    }
    while (synced < chat.rev) {
      final page = await _api.since(chat.id, synced);
      var next = page.messages.isEmpty ? synced : page.messages.last.rev;
      if (!page.more && next < chat.rev) {
        next = chat.rev;
      }
      await _store.saveMessages(chat.id, page.messages, syncedRev: next);
      if (next <= synced) {
        return;
      }
      synced = next;
    }
  }

  Future<void> older(String chatId) async {
    final first = await _store.firstSeq(chatId);
    if (first == null) {
      return;
    }
    final page = await _api.before(chatId, first);
    await _store.saveMessages(chatId, page.messages, older: page.more);
  }

  Future<void> read(Chat chat) async {
    if (chat.unread == 0) {
      return;
    }
    await _store.markRead(chat.id, chat.seq);
    try {
      await _api.read(chat.id, chat.seq);
    } on ApiException catch (error) {
      if (!offline(error)) {
        rethrow;
      }
    }
  }
}
