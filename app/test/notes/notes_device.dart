import 'package:drift/drift.dart';
import 'package:drift/native.dart';
import 'package:lodestar/core/api/api_client.dart';
import 'package:lodestar/core/api/api_transport.dart';
import 'package:lodestar/core/api/auth_api.dart';
import 'package:lodestar/core/api/server_address.dart';
import 'package:lodestar/core/api/token_pair.dart';
import 'package:lodestar/features/notes/note.dart';
import 'package:lodestar/features/notes/note_meta.dart';
import 'package:lodestar/features/notes/notes_api.dart';
import 'package:lodestar/features/notes/notes_database.dart';
import 'package:lodestar/features/notes/notes_store.dart';
import 'package:lodestar/features/notes/notes_sync.dart';

import '../support/fake_server.dart';
import '../support/memory_token_store.dart';

class NotesDevice {
  NotesDevice(FakeServer server)
    : db = NotesDatabase(
        DatabaseConnection(
          NativeDatabase.memory(),
          closeStreamsSynchronously: true,
        ),
      ) {
    final transport = ApiTransport(server.client);
    store = NotesStore(db, flushDelay: Duration.zero);
    sync = NotesSync(
      NotesApi(
        ApiClient(
          transport: transport,
          auth: AuthApi(transport),
          address: ServerAddress.tryParse('campus.example.edu', 'kaznu')!,
          tokens: MemoryTokenStore(
            const TokenPair(accessToken: 'access-1', refreshToken: 'r'),
          ),
          onExpired: () async {},
        ),
      ),
      store,
      compactAfter: 3,
    );
  }

  final NotesDatabase db;
  late final NotesStore store;
  late final NotesSync sync;

  Future<String> create(String path) => store.create(NoteMeta(path: path));

  Future<void> type(String id, String text) async {
    final doc = await store.open(id);
    store.edit(doc, text);
    await store.flush(id);
  }

  Future<String> text(String id) async => (await store.open(id)).text;

  Future<List<Note>> notes() => store.watchAll().first;

  Future<void> close() async {
    await store.clear();
    await db.close();
  }
}
