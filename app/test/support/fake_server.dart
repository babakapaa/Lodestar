import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:http/testing.dart';

class FakeServer {
  final List<http.Request> requests = [];
  final Map<String, Object?> profile = {
    'userId': 'u1',
    'displayName': '',
    'about': '',
    'contacts': <Object?>[],
  };
  final List<Map<String, Object?>> sessions = [
    {
      'id': 's1',
      'deviceName': 'Android',
      'createdAt': '2026-11-01T10:00:00Z',
      'lastUsedAt': '2026-11-12T10:00:00Z',
      'current': true,
    },
    {
      'id': 's2',
      'deviceName': 'Windows',
      'createdAt': '2026-11-02T10:00:00Z',
      'lastUsedAt': '2026-11-10T10:00:00Z',
      'current': false,
    },
  ];
  String validAccess = 'access-1';
  int issued = 1;
  bool acceptRefresh = true;

  MockClient get client => MockClient(_handle);

  Future<http.Response> _handle(http.Request request) async {
    requests.add(request);
    final path = request.url.path;
    final route = '${request.method} $path';
    if (path.startsWith('/v1/auth/')) {
      return _auth(request, route);
    }
    if (request.headers['Authorization'] != 'Bearer $validAccess') {
      return _json(401, {'code': 'unauthorized'});
    }
    return switch (route) {
      'GET /v1/me' => _json(200, {
        'id': 'u1',
        'tenantId': 't1',
        'email': 'aigerim@kaznu.kz',
        'displayName': 'Айгерим Сапарова',
        'createdAt': '2026-11-01T10:00:00Z',
      }),
      'GET /v1/me/profile' => _json(200, profile),
      'PUT /v1/me/profile' => _putProfile(request),
      'GET /v1/sessions' => _json(200, sessions),
      'DELETE /v1/sessions/s2' => _revoke(),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  http.Response _auth(http.Request request, String route) {
    switch (route) {
      case 'GET /v1/auth/kaznu/providers':
        return _json(200, [
          {'key': 'microsoft', 'displayName': 'Microsoft'},
        ]);
      case 'POST /v1/auth/token':
        final body = jsonDecode(request.body) as Map<String, Object?>;
        final rejected = body['grantType'] == 'refresh_token'
            ? !acceptRefresh
            : body['code'] != 'one-time-code';
        if (rejected) {
          return _json(400, {'code': 'invalid_grant'});
        }
        issued++;
        validAccess = 'access-$issued';
        return _json(200, {
          'accessToken': validAccess,
          'refreshToken': 'refresh-$issued',
          'tokenType': 'Bearer',
          'expiresIn': 900,
        });
      case 'POST /v1/auth/logout':
        return http.Response('', 204);
    }
    if (route.endsWith('/providers')) {
      return _json(404, {'code': 'tenant_not_found'});
    }
    return _json(404, {'code': 'not_found'});
  }

  http.Response _putProfile(http.Request request) {
    profile.addAll(jsonDecode(request.body) as Map<String, Object?>);
    return _json(200, profile);
  }

  http.Response _revoke() {
    sessions.removeWhere((session) => session['id'] == 's2');
    return http.Response('', 204);
  }

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
