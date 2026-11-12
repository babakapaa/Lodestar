import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeSources {
  final Map<String, Map<String, Object?>> sources = {};

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('GET', ['sources']) => _json(200, sources.values.toList()),
      ('POST', ['sources']) => _save(request),
      ('DELETE', ['sources', final id]) => _delete(id),
      ('POST', ['lookup']) => _lookup(request),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  http.Response _save(http.Request request) {
    final body = _body(request);
    final id = body['id']! as String;
    final created = !sources.containsKey(id);
    sources[id] = {...body, 'bytes': 100, 'createdAt': '2026-11-12T10:00:00Z'};
    return _json(created ? 201 : 200, sources[id]);
  }

  http.Response _delete(String id) {
    if (sources.remove(id) == null) {
      return _json(404, {'code': 'source_not_found'});
    }
    return http.Response('', 204);
  }

  http.Response _lookup(http.Request request) {
    if (_body(request)['doi'] != '10.1000/xyz') {
      return _json(422, {'code': 'source_lookup_failed'});
    }
    return _json(200, {
      'kind': 'article',
      'title': 'Строение клетки',
      'authors': ['Иванов, А. Б.'],
      'year': '2021',
      'container': 'Вестник биологии',
      'publisher': '',
      'city': '',
      'volume': '12',
      'issue': '3',
      'pages': '45-67',
      'url': '',
      'doi': '10.1000/xyz',
      'isbn': '',
      'accessedAt': '',
      'note': '',
    });
  }

  static Map<String, Object?> _body(http.Request request) =>
      jsonDecode(request.body) as Map<String, Object?>;

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
