import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeLti {
  final List<Map<String, Object?>> platforms = [];
  final Map<String, Map<String, Object?>> contexts = {
    'x1': {
      'id': 'x1',
      'title': 'Biology 101',
      'label': 'BIO',
      'platform': 'Moodle',
      'instructor': true,
      'link': null,
      'options': [
        {
          'kind': 'study.course',
          'id': 'c1',
          'title': 'Биология',
          'path': '/study/courses/c1',
        },
      ],
      'assignmentSource': 'lodestar',
      'sourceEditable': true,
      'canSync': false,
      'syncedAt': null,
      'syncError': null,
    },
  };
  int syncs = 0;
  final List<Object?> answers = [];
  int _next = 1;

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('GET', ['tool']) => _json(200, {
        'loginUrl': 'https://campus.example.edu/v1/lti/tenants/kaznu/login',
        'launchUrl': 'https://campus.example.edu/v1/lti/tenants/kaznu/launch',
        'jwksUrl': 'https://campus.example.edu/v1/lti/tenants/kaznu/jwks',
      }),
      ('GET', ['platforms']) => _json(200, platforms),
      ('POST', ['platforms']) => _create(request),
      ('DELETE', ['platforms', final id]) => _delete(id),
      ('GET', ['contexts', final id]) => _context(id),
      ('PUT', ['contexts', final id, 'link']) => _link(request, id),
      ('DELETE', ['contexts', final id, 'link']) => _unlink(id),
      ('PUT', ['contexts', final id, 'source']) => _source(request, id),
      ('POST', ['contexts', final id, 'sync']) => _sync(id),
      ('GET', ['links', 'study.course', final id]) => _json(200, [
        for (final c in contexts.values)
          if ((c['link'] as Map<String, Object?>?)?['id'] == id) c,
      ]),
      ('GET', ['deep-links', 'd1']) => _json(200, {
        'id': 'd1',
        'platform': 'Moodle',
        'acceptMultiple': false,
        'targets': contexts['x1']!['options'],
        'createdAt': '2026-11-13T10:00:00Z',
      }),
      ('POST', ['deep-links', 'd1']) => _answer(request),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  http.Response _answer(http.Request request) {
    answers.add(jsonDecode(request.body));
    return _json(200, {
      'formUrl': 'https://campus.example.edu/v1/lti/tenants/kaznu/deep-links/return?code=abc',
    });
  }

  http.Response _create(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    if (platforms.any((p) => p['issuer'] == body['issuer'])) {
      return _json(409, {'code': 'lti_platform_taken'});
    }
    final platform = {...body, 'id': 'p${_next++}'};
    platforms.add(platform);
    return _json(201, platform);
  }

  http.Response _delete(String id) {
    platforms.removeWhere((p) => p['id'] == id);
    return http.Response('', 204);
  }

  http.Response _context(String id) {
    final c = contexts[id];
    return c == null
        ? _json(404, {'code': 'lti_context_not_found'})
        : _json(200, c);
  }

  http.Response _link(http.Request request, String id) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final c = contexts[id]!;
    c['link'] = (c['options']! as List<Object?>)
        .cast<Map<String, Object?>>()
        .firstWhere((o) => o['id'] == body['id']);
    c['canSync'] = true;
    return _json(200, c);
  }

  http.Response _unlink(String id) {
    final c = contexts[id]!;
    c['link'] = null;
    c['canSync'] = false;
    return _json(200, c);
  }

  http.Response _source(http.Request request, String id) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final c = contexts[id]!;
    c['assignmentSource'] = body['source'];
    return _json(200, c);
  }

  http.Response _sync(String id) {
    syncs++;
    contexts[id]!['syncedAt'] = '2026-11-13T10:00:00Z';
    return _json(200, {'members': 3, 'enrolled': 2, 'works': 1, 'scores': 1});
  }

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
