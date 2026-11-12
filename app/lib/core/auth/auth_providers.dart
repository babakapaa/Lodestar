import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:http/http.dart' as http;

import '../api/api_client.dart';
import '../api/api_transport.dart';
import '../api/auth_api.dart';
import '../api/connection_store.dart';
import '../api/token_store.dart';
import '../settings/settings_providers.dart';
import 'app_redirect.dart';
import 'auth_controller.dart';
import 'auth_state.dart';
import 'browser_authenticator.dart';
import 'device_name.dart';
import 'web_authenticator.dart';

final httpClientProvider = Provider<http.Client>((ref) {
  final client = http.Client();
  ref.onDispose(client.close);
  return client;
});

final apiTransportProvider = Provider<ApiTransport>(
  (ref) => ApiTransport(ref.watch(httpClientProvider)),
);

final authApiProvider = Provider<AuthApi>(
  (ref) => AuthApi(ref.watch(apiTransportProvider)),
);

final tokenStoreProvider = Provider<TokenStore>(
  (ref) => throw StateError('tokenStoreProvider must be overridden at startup'),
);

final connectionStoreProvider = Provider<ConnectionStore>(
  (ref) => ConnectionStore(ref.watch(sharedPreferencesProvider)),
);

final webAuthenticatorProvider = Provider<WebAuthenticator>(
  (ref) => const BrowserAuthenticator(),
);

final appRedirectProvider = Provider<AppRedirect>(
  (ref) => AppRedirect.current(),
);

final deviceNameProvider = Provider<String>((ref) => DeviceName.current());

final authProvider = NotifierProvider<AuthController, AuthState>(
  AuthController.new,
);

final apiClientProvider = Provider<ApiClient>((ref) {
  final address = ref.watch(authProvider.select((state) => state.address));
  if (address == null) {
    throw StateError('apiClientProvider needs a connected server');
  }
  return ApiClient(
    transport: ref.watch(apiTransportProvider),
    auth: ref.watch(authApiProvider),
    address: address,
    tokens: ref.watch(tokenStoreProvider),
    onExpired: () => ref.read(authProvider.notifier).expire(),
  );
});
