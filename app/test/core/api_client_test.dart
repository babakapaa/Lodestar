import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/api/api_client.dart';
import 'package:lodestar/core/api/api_exception.dart';
import 'package:lodestar/core/api/api_transport.dart';
import 'package:lodestar/core/api/auth_api.dart';
import 'package:lodestar/core/api/server_address.dart';
import 'package:lodestar/core/api/token_pair.dart';

import '../support/fake_server.dart';
import '../support/memory_token_store.dart';

void main() {
  final address = ServerAddress.tryParse('campus.example.edu', 'kaznu')!;

  ApiClient clientFor(
    FakeServer server,
    MemoryTokenStore tokens, {
    void Function()? onExpired,
  }) {
    final transport = ApiTransport(server.client);
    return ApiClient(
      transport: transport,
      auth: AuthApi(transport),
      address: address,
      tokens: tokens,
      onExpired: () async => onExpired?.call(),
    );
  }

  test(
    'an expired access token is refreshed once for parallel calls',
    () async {
      final server = FakeServer()..validAccess = 'access-0';
      final tokens = MemoryTokenStore(
        const TokenPair(accessToken: 'stale', refreshToken: 'refresh-0'),
      );
      final client = clientFor(server, tokens);

      await Future.wait([client.get('/v1/me'), client.get('/v1/sessions')]);

      final refreshes = server.requests.where(
        (request) => request.url.path == '/v1/auth/token',
      );
      expect(refreshes, hasLength(1));
      expect(tokens.current!.accessToken, server.validAccess);
    },
  );

  test('a rejected refresh signs the user out', () async {
    final server = FakeServer()..acceptRefresh = false;
    final tokens = MemoryTokenStore(
      const TokenPair(accessToken: 'stale', refreshToken: 'refresh-0'),
    );
    var expired = false;
    final client = clientFor(server, tokens, onExpired: () => expired = true);

    await expectLater(
      client.get('/v1/me/profile'),
      throwsA(isA<ApiException>()),
    );
    expect(expired, isTrue);
  });

  test('server errors keep their code', () async {
    final server = FakeServer();
    final client = clientFor(
      server,
      MemoryTokenStore(
        const TokenPair(accessToken: 'access-1', refreshToken: 'refresh-1'),
      ),
    );

    await expectLater(
      client.get('/v1/unknown'),
      throwsA(
        isA<ApiException>()
            .having((error) => error.statusCode, 'status', 404)
            .having((error) => error.code, 'code', 'not_found'),
      ),
    );
  });
}
