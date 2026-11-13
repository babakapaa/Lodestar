import 'package:drift/drift.dart';

import 'chat_rows.dart';
import 'chat_values.dart';
import 'message_rows.dart';
import 'outbox_rows.dart';

part 'chats_database.g.dart';

@DriftDatabase(tables: [ChatRows, MessageRows, OutboxRows, ChatValues])
class ChatsDatabase extends _$ChatsDatabase {
  ChatsDatabase(super.executor);

  @override
  int get schemaVersion => 1;

  @override
  MigrationStrategy get migration => MigrationStrategy(
    onCreate: (m) async {
      await m.createAll();
      await customStatement(
        'CREATE INDEX message_rows_chat ON message_rows (chat_id, seq)',
      );
    },
  );

  Future<void> clear() => transaction(() async {
    await delete(outboxRows).go();
    await delete(messageRows).go();
    await delete(chatRows).go();
    await delete(chatValues).go();
  });
}
