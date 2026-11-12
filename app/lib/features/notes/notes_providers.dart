import 'package:drift_flutter/drift_flutter.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../../core/offline/loaded.dart';
import '../../core/offline/offline_providers.dart';
import 'note.dart';
import 'note_member.dart';
import 'note_version.dart';
import 'notes_api.dart';
import 'notes_database.dart';
import 'notes_store.dart';
import 'notes_sync.dart';
import 'notes_sync_controller.dart';
import 'notes_sync_status.dart';
import 'source.dart';
import 'sources_api.dart';

final notesDatabaseProvider = Provider<NotesDatabase>((ref) {
  final db = NotesDatabase(
    driftDatabase(
      name: 'lodestar_notes',
      web: DriftWebOptions(
        sqlite3Wasm: Uri.parse('sqlite3.wasm'),
        driftWorker: Uri.parse('drift_worker.js'),
      ),
    ),
  );
  ref.onDispose(db.close);
  return db;
});

final notesStoreProvider = Provider<NotesStore>((ref) {
  final store = NotesStore(ref.watch(notesDatabaseProvider));
  ref.onDispose(store.dispose);
  return store;
});

final notesApiProvider = Provider<NotesApi>(
  (ref) => NotesApi(ref.watch(apiClientProvider)),
);

final notesSyncProvider = Provider<NotesSync>(
  (ref) =>
      NotesSync(ref.watch(notesApiProvider), ref.watch(notesStoreProvider)),
);

final notesSyncStatusProvider =
    NotifierProvider<NotesSyncController, NotesSyncStatus>(
      NotesSyncController.new,
    );

final notesProvider = StreamProvider<List<Note>>(
  (ref) => ref.watch(notesStoreProvider).watchAll(),
);

final noteProvider = StreamProvider.autoDispose.family<Note?, String>(
  (ref, id) => ref.watch(notesStoreProvider).watch(id),
);

final notesPendingProvider = StreamProvider<bool>(
  (ref) => ref.watch(notesStoreProvider).watchPending(),
);

final noteMembersProvider = FutureProvider.autoDispose
    .family<List<NoteMember>, String>(
      (ref, id) => ref.watch(notesApiProvider).members(id),
    );

final noteVersionsProvider = FutureProvider.autoDispose
    .family<List<NoteVersion>, String>(
      (ref, id) => ref.watch(notesApiProvider).versions(id),
    );

final noteVersionBodyProvider = FutureProvider.autoDispose
    .family<String, (String, String)>(
      (ref, key) => ref.watch(notesApiProvider).versionBody(key.$1, key.$2),
    );

final sourcesApiProvider = Provider<SourcesApi>(
  (ref) =>
      SourcesApi(ref.watch(apiClientProvider), ref.watch(cachedReaderProvider)),
);

final sourcesProvider = FutureProvider.autoDispose<Loaded<List<Source>>>(
  (ref) => ref.watch(sourcesApiProvider).list(),
);
