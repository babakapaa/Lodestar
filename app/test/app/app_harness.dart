import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/app/lodestar_app.dart';
import 'package:lodestar/core/api/token_pair.dart';
import 'package:lodestar/core/auth/auth_providers.dart';
import 'package:lodestar/core/settings/settings_providers.dart';
import 'package:material_ui/material_ui.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../support/fake_authenticator.dart';
import '../support/fake_server.dart';
import '../support/memory_token_store.dart';

const signedInStorage = <String, Object>{
  'settings.language': 'ru',
  'connection.server': 'https://campus.example.edu',
  'connection.tenant': 'kaznu',
};

Future<FakeServer> pumpApp(
  WidgetTester tester, {
  required Size size,
  Map<String, Object> stored = signedInStorage,
  bool signedIn = true,
  FakeAuthenticator? authenticator,
}) async {
  tester.view.physicalSize = size;
  tester.view.devicePixelRatio = 1;
  addTearDown(tester.view.reset);
  SharedPreferences.setMockInitialValues(stored);
  final prefs = await SharedPreferences.getInstance();
  final server = FakeServer();
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
        httpClientProvider.overrideWithValue(server.client),
        webAuthenticatorProvider.overrideWithValue(
          authenticator ?? FakeAuthenticator(),
        ),
      ],
      child: const LodestarApp(),
    ),
  );
  await tester.pumpAndSettle();
  return server;
}
