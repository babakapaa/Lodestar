import '../auth/identity_provider.dart';
import 'api_transport.dart';
import 'server_address.dart';
import 'token_pair.dart';

class AuthApi {
  const AuthApi(this._transport);

  final ApiTransport _transport;

  Future<List<IdentityProvider>> providers(ServerAddress address) async {
    final body = await _transport.send(
      'GET',
      address.resolve('/v1/auth/${address.tenant}/providers'),
    );
    return [
      for (final item in body! as List<Object?>)
        IdentityProvider.fromJson(item! as Map<String, Object?>),
    ];
  }

  Uri loginUri(
    ServerAddress address,
    IdentityProvider provider, {
    required Uri redirect,
    required String challenge,
    required String device,
  }) {
    return address.resolve(
      '/v1/auth/${address.tenant}/${Uri.encodeComponent(provider.key)}/login',
      {
        'redirect_uri': redirect.toString(),
        'code_challenge': challenge,
        'code_challenge_method': 'S256',
        'device': device,
      },
    );
  }

  Future<TokenPair> exchangeCode(
    ServerAddress address, {
    required String code,
    required String verifier,
    required Uri redirect,
  }) {
    return _issue(address, {
      'grantType': 'authorization_code',
      'code': code,
      'codeVerifier': verifier,
      'redirectUri': redirect.toString(),
    });
  }

  Future<TokenPair> refresh(ServerAddress address, TokenPair tokens) {
    return _issue(address, {
      'grantType': 'refresh_token',
      'refreshToken': tokens.refreshToken,
    });
  }

  Future<void> logout(ServerAddress address, TokenPair tokens) async {
    await _transport.send(
      'POST',
      address.resolve('/v1/auth/logout'),
      accessToken: tokens.accessToken,
    );
  }

  Future<TokenPair> _issue(
    ServerAddress address,
    Map<String, String> grant,
  ) async {
    final body = await _transport.send(
      'POST',
      address.resolve('/v1/auth/token'),
      body: grant,
    );
    return TokenPair.fromJson(body! as Map<String, Object?>);
  }
}
