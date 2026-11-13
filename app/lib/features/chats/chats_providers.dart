import 'package:drift_flutter/drift_flutter.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../../core/offline/offline_providers.dart';
import 'blocked_user.dart';
import 'chat.dart';
import 'chat_details.dart';
import 'chat_entry.dart';
import 'chat_report.dart';
import 'chats_api.dart';
import 'chats_database.dart';
import 'chats_store.dart';
import 'chats_sync.dart';
import 'chats_sync_controller.dart';
import 'chats_sync_status.dart';

final chatsDatabaseProvider = Provider<ChatsDatabase>((ref) {
  final db = ChatsDatabase(
    driftDatabase(
      name: 'lodestar_chats',
      web: DriftWebOptions(
        sqlite3Wasm: Uri.parse('sqlite3.wasm'),
        driftWorker: Uri.parse('drift_worker.js'),
      ),
    ),
  );
  ref.onDispose(db.close);
  return db;
});

final chatsStoreProvider = Provider<ChatsStore>((ref) {
  final store = ChatsStore(ref.watch(chatsDatabaseProvider));
  ref.onDispose(store.dispose);
  return store;
});

final chatsApiProvider = Provider<ChatsApi>(
  (ref) =>
      ChatsApi(ref.watch(apiClientProvider), ref.watch(cachedReaderProvider)),
);

final chatsSyncProvider = Provider<ChatsSync>(
  (ref) =>
      ChatsSync(ref.watch(chatsApiProvider), ref.watch(chatsStoreProvider)),
);

final chatsSyncStatusProvider =
    NotifierProvider<ChatsSyncController, ChatsSyncStatus>(
      ChatsSyncController.new,
    );

final chatsProvider = StreamProvider<List<Chat>>(
  (ref) => ref.watch(chatsStoreProvider).watchChats(),
);

final chatProvider = StreamProvider.autoDispose.family<Chat?, String>(
  (ref, id) => ref.watch(chatsStoreProvider).watchChat(id),
);

final chatEntriesProvider = StreamProvider.autoDispose
    .family<List<ChatEntry>, String>(
      (ref, id) => ref.watch(chatsStoreProvider).watchEntries(id),
    );

final chatOlderProvider = StreamProvider.autoDispose.family<bool, String>(
  (ref, id) => ref.watch(chatsStoreProvider).watchOlder(id),
);

final chatDetailsProvider = FutureProvider.autoDispose
    .family<ChatDetails, String>(
      (ref, id) => ref.watch(chatsApiProvider).get(id),
    );

final myIdProvider = StreamProvider<String?>(
  (ref) => ref.watch(chatsStoreProvider).watchMe(),
);

final chatBlocksProvider = FutureProvider.autoDispose<List<BlockedUser>>(
  (ref) => ref.watch(chatsApiProvider).blocks(),
);

final chatReportsProvider = FutureProvider.autoDispose<List<ChatReport>>(
  (ref) => ref.watch(chatsApiProvider).reports(),
);
