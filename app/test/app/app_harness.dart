import 'package:drift/drift.dart' show DatabaseConnection;
import 'package:drift/native.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/app/lodestar_app.dart';
import 'package:lodestar/core/api/token_pair.dart';
import 'package:lodestar/core/auth/auth_providers.dart';
import 'package:lodestar/core/live/live_providers.dart';
import 'package:lodestar/core/settings/settings_providers.dart';
import 'package:lodestar/features/notes/notes_database.dart';
import 'package:lodestar/features/notes/notes_providers.dart';
import 'package:lodestar/features/privacy/privacy_providers.dart';
import 'package:material_ui/material_ui.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../support/fake_authenticator.dart';
import '../support/fake_export_saver.dart';
import '../support/fake_live_socket.dart';
import '../support/fake_server.dart';
import '../support/memory_token_store.dart';

const signedInStorage = <String, Object>{
  'settings.language': 'ru',
  'connection.server': 'https://campus.example.edu',
  'connection.tenant': 'kaznu',
};

class AppHarness {
  const AppHarness(this.server, this.live, this.saver);

  final FakeServer server;
  final FakeLiveSocket live;
  final FakeExportSaver saver;
}

Future<AppHarness> pumpApp(
  WidgetTester tester, {
  required Size size,
  Map<String, Object> stored = signedInStorage,
  bool signedIn = true,
  FakeAuthenticator? authenticator,
  FakeServer? server,
}) async {
  tester.view.physicalSize = size;
  tester.view.devicePixelRatio = 1;
  addTearDown(tester.view.reset);
  SharedPreferences.setMockInitialValues(stored);
  final prefs = await SharedPreferences.getInstance();
  final fake = server ?? FakeServer();
  final live = FakeLiveSocket();
  final saver = FakeExportSaver();
  final tokens = MemoryTokenStore(
    signedIn
        ? const TokenPair(accessToken: 'access-1', refreshToken: 'refresh-1')
        : null,
  );
  await tester.pumpWidget(
    ProviderScope(
      overrides: [
        sharedPreferencesProvider.overrideWithValue(prefs),
        tokenStoreProvider.overrideWithValue(tokens),
        httpClientProvider.overrideWithValue(fake.client),
        liveSocketProvider.overrideWithValue(live),
        exportSaverProvider.overrideWithValue(saver),
        notesDatabaseProvider.overrideWith((ref) {
          final db = NotesDatabase(
            DatabaseConnection(
              NativeDatabase.memory(),
              closeStreamsSynchronously: true,
            ),
          );
          ref.onDispose(db.close);
          return db;
        }),
        webAuthenticatorProvider.overrideWithValue(
          authenticator ?? FakeAuthenticator(),
        ),
      ],
      child: const LodestarApp(),
    ),
  );
  await tester.pumpAndSettle();
  return AppHarness(fake, live, saver);
}
