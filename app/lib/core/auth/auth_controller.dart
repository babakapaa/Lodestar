import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../api/api_exception.dart';
import '../api/server_address.dart';
import '../offline/offline_providers.dart';
import 'auth_providers.dart';
import 'auth_state.dart';
import 'identity_provider.dart';
import 'pkce.dart';

class AuthController extends Notifier<AuthState> {
  @override
  AuthState build() {
    final address = ref.read(connectionStoreProvider).load();
    if (address == null) {
      return const Disconnected();
    }
    if (ref.read(tokenStoreProvider).current == null) {
      return SignedOut(address);
    }
    return SignedIn(address);
  }

  Future<void> connect(ServerAddress address) async {
    await ref.read(authApiProvider).providers(address);
    await ref.read(connectionStoreProvider).save(address);
    state = SignedOut(address);
  }

  Future<void> disconnect() async {
    await ref.read(tokenStoreProvider).clear();
    await ref.read(responseCacheProvider).clear();
    await ref.read(connectionStoreProvider).clear();
    state = const Disconnected();
  }

  Future<bool> signIn(IdentityProvider provider) async {
    final address = state.address;
    if (address == null) {
      return false;
    }
    final api = ref.read(authApiProvider);
    final redirect = ref.read(appRedirectProvider);
    final pkce = Pkce.generate();
    final result = await ref
        .read(webAuthenticatorProvider)
        .authenticate(
          api.loginUri(
            address,
            provider,
            redirect: redirect.uri,
            challenge: pkce.challenge,
            device: ref.read(deviceNameProvider),
          ),
          redirect.callbackScheme,
        );
    if (result == null) {
      return false;
    }
    final code = result.queryParameters['code'];
    if (code == null || code.isEmpty) {
      throw ApiException(
        400,
        result.queryParameters['error'] ?? ApiException.loginFailed,
      );
    }
    final tokens = await api.exchangeCode(
      address,
      code: code,
      verifier: pkce.verifier,
      redirect: redirect.uri,
    );
    await ref.read(tokenStoreProvider).save(tokens);
    state = SignedIn(address);
    return true;
  }

  Future<void> signOut() async {
    final address = state.address;
    final tokens = ref.read(tokenStoreProvider).current;
    if (address != null && tokens != null) {
      try {
        await ref.read(authApiProvider).logout(address, tokens);
      } on ApiException {
        await expire();
        return;
      }
    }
    await expire();
  }

  Future<void> expire() async {
    final address = state.address;
    await ref.read(tokenStoreProvider).clear();
    await ref.read(responseCacheProvider).clear();
    if (address != null) {
      state = SignedOut(address);
    }
  }
}
