import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:http/testing.dart';

import 'fake_directory.dart';
import 'fake_navigator.dart';
import 'fake_notes.dart';
import 'fake_sources.dart';
import 'fake_study.dart';

class FakeServer {
  final List<http.Request> requests = [];
  final FakeDirectory directory = FakeDirectory();
  final FakeStudy study = FakeStudy();
  final FakeNotes notes = FakeNotes();
  final FakeSources sources = FakeSources();
  final FakeNavigator navigator = FakeNavigator();
  bool offline = false;
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
  final List<Map<String, Object?>> notifications = [
    {
      'id': 'n1',
      'type': 'study.deadline',
      'params': {'course': 'Биология'},
      'link': '',
      'createdAt': '2026-11-12T09:00:00Z',
      'readAt': null,
    },
  ];
  final Map<String, Map<String, Object?>> channels = {
    'study.deadline': {
      'type': 'study.deadline',
      'required': false,
      'inApp': true,
      'push': true,
      'email': false,
    },
  };
  final Map<String, Object?> preferences = {};
  Map<String, Object?>? deletion;
  List<String> permissions = [];
  final Map<String, Object?> settings = {};
  final List<Map<String, Object?>> modules = [
    {'name': 'study', 'version': '0.1.0', 'optional': true, 'enabled': true},
  ];
  bool selfDeletion = true;
  String validAccess = 'access-1';
  int issued = 1;
  bool acceptRefresh = true;

  MockClient get client => MockClient(_handle);

  Future<http.Response> _handle(http.Request request) async {
    if (offline) {
      throw http.ClientException('offline', request.url);
    }
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
      'GET /v1/notifications' => _json(200, notifications),
      'GET /v1/notifications/unread' => _json(200, {
        'unread': notifications.where((n) => n['readAt'] == null).length,
      }),
      'POST /v1/notifications/read-all' => _readAll(),
      'POST /v1/notifications/n1/read' => _readAll(),
      'GET /v1/notification-types' => _json(200, [
        {
          'key': 'study.deadline',
          'required': false,
          'names': {'ru': 'Сроки сдачи', 'en': 'Deadlines'},
          'texts': {
            'ru': {'title': 'Скоро срок: {course}', 'body': 'Осталось 2 дня'},
            'en': {'title': 'Deadline soon: {course}', 'body': '2 days left'},
          },
        },
      ]),
      'GET /v1/me/notification-channels' => _json(
        200,
        channels.values.toList(),
      ),
      'PUT /v1/me/notification-channels/study.deadline' => _setChannels(
        request,
      ),
      'GET /v1/me/preferences' => _json(200, [
        for (final key in [
          'core.theme',
          'core.appearance',
          'core.language',
          'core.quiet_hours',
          'core.quiet_from',
          'core.quiet_to',
          'navigator.journal',
        ])
          {
            'key': key,
            'value': preferences[key] ?? _defaults[key],
            'set': preferences.containsKey(key),
          },
      ]),
      'GET /v1/me/permissions' => _json(200, permissions),
      'GET /v1/settings/schema' => _json(200, [
        {
          'key': 'core.default_language',
          'type': 'choice',
          'options': ['system', 'kk', 'ru', 'en'],
          'default': 'system',
        },
        {'key': 'core.self_deletion', 'type': 'bool', 'default': true},
        {
          'key': 'core.deletion_grace_days',
          'type': 'integer',
          'min': 0,
          'max': 90,
          'default': 14,
        },
      ]),
      'GET /v1/settings' => _json(200, [
        for (final key in [
          'core.default_language',
          'core.self_deletion',
          'core.deletion_grace_days',
        ])
          {
            'key': key,
            'value': settings[key] ?? _settingDefaults[key],
            'set': settings.containsKey(key),
          },
      ]),
      'GET /v1/modules' => _json(200, modules),
      'PUT /v1/modules/study' => _setModule(request),
      'GET /v1/audit' => _json(200, [
        {
          'id': 3,
          'actorId': 'u1',
          'action': 'core.setting.changed',
          'targetKind': 'setting',
          'targetId': 'core.self_deletion',
          'details': {'key': 'core.self_deletion'},
          'at': '2026-11-12T09:30:00Z',
        },
      ]),
      'GET /v1/users/u1/profile' => _json(200, {
        'userId': 'u1',
        'displayName': 'Дана Ахметова',
        'about': '',
        'contacts': <Object?>[],
      }),
      _ when route.startsWith('PUT /v1/settings/') => _setSetting(request),
      'GET /v1/me/export' => http.Response.bytes([80, 75, 3, 4], 200),
      'GET /v1/me/deletion' =>
        deletion == null
            ? _json(404, {'code': 'deletion_not_requested'})
            : _json(200, deletion),
      'POST /v1/me/deletion' => _requestDeletion(),
      'DELETE /v1/me/deletion' => _cancelDeletion(),
      _ when route.startsWith('PUT /v1/me/preferences/') => _setPreference(
        request,
      ),
      _ when path.startsWith('/v1/study/') => study.handle(request),
      _ when path.startsWith('/v1/navigator/') => navigator.handle(request),
      _
          when path.startsWith('/v1/notes/sources') ||
              path == '/v1/notes/lookup' =>
        sources.handle(request),
      _ when path == '/v1/notes' || path.startsWith('/v1/notes/') =>
        notes.handle(request),
      _ => directory.handle(request),
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

  static const _defaults = <String, Object?>{
    'core.theme': 'tenant',
    'core.appearance': 'system',
    'core.language': 'system',
    'core.quiet_hours': false,
    'core.quiet_from': 1320,
    'core.quiet_to': 420,
    'navigator.journal': true,
  };

  static const _settingDefaults = <String, Object?>{
    'core.default_language': 'system',
    'core.self_deletion': true,
    'core.deletion_grace_days': 14,
  };

  http.Response _setSetting(http.Request request) {
    final key = Uri.decodeComponent(request.url.pathSegments.last);
    settings[key] = (jsonDecode(request.body) as Map<String, Object?>)['value'];
    return http.Response('', 204);
  }

  http.Response _setModule(http.Request request) {
    modules.first['enabled'] =
        (jsonDecode(request.body) as Map<String, Object?>)['enabled'];
    return http.Response('', 204);
  }

  http.Response _readAll() {
    for (final n in notifications) {
      n['readAt'] ??= '2026-11-12T10:00:00Z';
    }
    return http.Response('', 204);
  }

  http.Response _setChannels(http.Request request) {
    channels['study.deadline']!.addAll(
      jsonDecode(request.body) as Map<String, Object?>,
    );
    return http.Response('', 204);
  }

  http.Response _setPreference(http.Request request) {
    final key = Uri.decodeComponent(request.url.pathSegments.last);
    preferences[key] =
        (jsonDecode(request.body) as Map<String, Object?>)['value'];
    return http.Response('', 204);
  }

  http.Response _requestDeletion() {
    if (!selfDeletion) {
      return _json(403, {'code': 'deletion_by_university'});
    }
    deletion ??= {
      'requestedAt': '2026-11-12T10:00:00Z',
      'dueAt': '2026-11-26T10:00:00Z',
    };
    return _json(200, deletion);
  }

  http.Response _cancelDeletion() {
    deletion = null;
    return http.Response('', 204);
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
