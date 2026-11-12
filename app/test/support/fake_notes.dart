import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeNotes {
  final Map<String, Map<String, Object?>> notes = {};
  final Map<String, List<Map<String, Object?>>> updates = {};
  final Map<String, Map<String, Object?>> snapshots = {};
  final Map<String, List<Map<String, Object?>>> members = {};
  int _seq = 0;

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('GET', []) => _json(200, notes.values.toList()),
      ('POST', []) => _create(request),
      ('PUT', [final id, 'meta']) => _meta(id, request),
      ('DELETE', [final id]) => _setDeleted(id, '2026-11-12T12:00:00Z'),
      ('POST', [final id, 'restore']) => _setDeleted(id, null),
      ('GET', [final id, 'updates']) => _changes(id, request),
      ('POST', [final id, 'updates']) => _push(id, request),
      ('POST', [final id, 'versions']) => _compact(id, request),
      ('GET', [final id, 'members']) => _json(200, members[id] ?? const []),
      ('PUT', [final id, 'members', final user]) => _share(id, user, request),
      ('DELETE', [final id, 'members', final user]) => _unshare(id, user),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  http.Response _create(http.Request request) {
    final body = _body(request);
    final id = body['id']! as String;
    if (notes.containsKey(id)) {
      return _json(200, notes[id]);
    }
    notes[id] = {
      'id': id,
      'ownerId': 'u1',
      'role': 'owner',
      'meta': body['meta'],
      'seq': 0,
      'bytes': 0,
      'createdAt': '2026-11-12T10:00:00Z',
      'changedAt': '2026-11-12T10:00:00Z',
      'deletedAt': null,
    };
    updates[id] = [];
    return _json(201, notes[id]);
  }

  http.Response _meta(String id, http.Request request) {
    final note = notes[id];
    if (note == null) {
      return _json(404, {'code': 'not_found'});
    }
    note['meta'] = _body(request);
    return _json(200, note);
  }

  http.Response _setDeleted(String id, String? at) {
    final note = notes[id];
    if (note == null) {
      return _json(404, {'code': 'not_found'});
    }
    note['deletedAt'] = at;
    return http.Response('', 204);
  }

  http.Response _changes(String id, http.Request request) {
    final after = int.parse(request.url.queryParameters['after'] ?? '0');
    final snapshot = snapshots[id];
    final useSnapshot = snapshot != null && after < (snapshot['seq']! as int);
    final from = useSnapshot ? snapshot['seq']! as int : after;
    return _json(200, {
      'seq': notes[id]!['seq'],
      'snapshot': useSnapshot ? snapshot : null,
      'updates': [
        for (final update in updates[id]!)
          if ((update['seq']! as int) > from) update,
      ],
    });
  }

  http.Response _push(String id, http.Request request) {
    final note = notes[id];
    if (note == null) {
      return _json(404, {'code': 'not_found'});
    }
    final body = _body(request);
    for (final update in updates[id]!) {
      if (update['key'] == body['key']) {
        return _json(200, {'seq': update['seq']});
      }
    }
    final seq = ++_seq;
    updates[id]!.add({
      'seq': seq,
      'key': body['key'],
      'authorId': 'u1',
      'data': body['data'],
      'createdAt': '2026-11-12T10:00:00Z',
    });
    note['seq'] = seq;
    return _json(201, {'seq': seq});
  }

  http.Response _compact(String id, http.Request request) {
    final body = _body(request);
    final seq = body['seq']! as int;
    snapshots[id] = {'seq': seq, 'data': body['data']};
    updates[id]!.removeWhere((update) => (update['seq']! as int) <= seq);
    return _json(201, {
      'id': 'v$seq',
      'seq': seq,
      'authorId': 'u1',
      'bytes': 0,
      'createdAt': '2026-11-12T10:00:00Z',
    });
  }

  http.Response _share(String id, String user, http.Request request) {
    final list = members.putIfAbsent(id, () => []);
    list
      ..removeWhere((member) => member['userId'] == user)
      ..add({
        'userId': user,
        'displayName': 'Ерлан Ким',
        'email': 'erlan@kaznu.kz',
        'canEdit': _body(request)['canEdit'],
        'addedAt': '2026-11-12T10:00:00Z',
      });
    return http.Response('', 204);
  }

  http.Response _unshare(String id, String user) {
    members[id]?.removeWhere((member) => member['userId'] == user);
    return http.Response('', 204);
  }

  static Map<String, Object?> _body(http.Request request) =>
      jsonDecode(request.body) as Map<String, Object?>;

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
