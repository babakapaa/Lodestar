import 'dart:convert';
import 'dart:typed_data';

import 'package:http/http.dart' as http;

class FakeFiles {
  static const host = 'storage.example';

  final Map<String, Map<String, Object?>> files = {};
  final Map<String, Uint8List> stored = {};
  int _next = 1;

  void add(String id, String name, List<int> bytes) {
    files[id] = _file(id, name, bytes.length);
    stored[id] = Uint8List.fromList(bytes);
  }

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('POST', []) => _create(jsonDecode(request.body) as Map<String, Object?>),
      ('POST', [final id, 'complete']) => _json(200, files[id]),
      ('GET', [final id, 'download']) => _json(200, {
        'url': 'https://$host/$id',
      }),
      _ => _json(404, {'code': 'file_not_found'}),
    };
  }

  http.Response storage(http.Request request) {
    final id = request.url.pathSegments.single;
    if (request.method == 'PUT') {
      stored[id] = request.bodyBytes;
      return http.Response('', 200);
    }
    final bytes = stored[id];
    return bytes == null
        ? http.Response('', 404)
        : http.Response.bytes(bytes, 200);
  }

  http.Response _create(Map<String, Object?> body) {
    final id = 'f${_next++}';
    files[id] = _file(id, body['name']! as String, body['size']! as int);
    return _json(201, {
      'file': files[id],
      'upload': {
        'url': 'https://$host/$id',
        'method': 'PUT',
        'headers': {'content-type': body['contentType']},
        'expiresAt': '2026-11-13T11:00:00Z',
      },
    });
  }

  static Map<String, Object?> _file(String id, String name, int size) => {
    'id': id,
    'ownerId': 'u1',
    'name': name,
    'contentType': 'application/pdf',
    'size': size,
    'status': 'uploaded',
    'createdAt': '2026-11-13T10:00:00Z',
    'uploadedAt': '2026-11-13T10:00:00Z',
  };

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
