import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeDirectory {
  final List<Map<String, Object?>> roles = [
    {
      'id': 'r-admin',
      'name': 'Администратор',
      'color': '#C0392B',
      'icon': '',
      'position': 2,
      'permissions': ['core.administrator'],
      'default': false,
    },
    {
      'id': 'r-teacher',
      'name': 'Преподаватель',
      'color': '',
      'icon': '',
      'position': 1,
      'permissions': ['core.structure.members'],
      'default': false,
    },
    {
      'id': 'r-everyone',
      'name': 'Все',
      'color': '',
      'icon': '',
      'position': 0,
      'permissions': <String>[],
      'default': true,
    },
  ];
  final List<Map<String, Object?>> assignments = [];
  final List<Map<String, Object?>> units = [
    {
      'id': 'f1',
      'parentId': null,
      'kind': 'faculty',
      'name': 'Факультет информатики',
      'code': null,
      'archivedAt': null,
      'createdAt': '2026-11-01T10:00:00Z',
    },
  ];
  final Map<String, Set<String>> members = {};
  final List<(String, int)> moves = [];
  int _next = 1;

  static const people = {'u2': ('Ерлан Ким', 'erlan@kaznu.kz')};

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(1).toList();
    final method = request.method;
    return switch ((method, segments)) {
      ('GET', ['permissions']) => _json(200, [
        'core.administrator',
        'core.roles.manage',
        'core.structure.manage',
        'core.structure.members',
      ]),
      ('GET', ['roles']) => _json(200, _sortedRoles()),
      ('POST', ['roles']) => _createRole(request),
      ('PUT', ['roles', final id]) => _updateRole(request, id),
      ('DELETE', ['roles', final id]) => _deleteRole(id),
      ('POST', ['roles', final id, 'move']) => _moveRole(request, id),
      ('GET', ['users', final id, 'roles']) => _json(200, [
        for (final item in assignments)
          if (item['userId'] == id) item,
      ]),
      ('GET', ['users', final id, 'profile']) when people.containsKey(id) =>
        _json(200, {
          'userId': id,
          'displayName': people[id]!.$1,
          'about': '',
          'contacts': <Object?>[],
        }),
      ('POST', ['role-assignments']) => _assign(request),
      ('DELETE', ['role-assignments', final id]) => _unassign(id),
      ('GET', ['search']) => _json(200, [
        for (final MapEntry(:key, :value) in people.entries)
          if (value.$1.toLowerCase().contains(
            request.url.queryParameters['q']!.toLowerCase(),
          ))
            {
              'kind': 'people',
              'id': key,
              'title': value.$1,
              'subtitle': value.$2,
              'link': '',
              'score': 1,
            },
      ]),
      ('GET', ['units']) => _json(200, [
        for (final unit in units)
          if (unit['parentId'] == request.url.queryParameters['parent'] &&
              (unit['archivedAt'] == null ||
                  request.url.queryParameters['archived'] == 'true'))
            unit,
      ]),
      ('POST', ['units']) => _createUnit(request),
      ('GET', ['units', final id]) => _json(200, {
        ..._unit(id),
        'path': _path(id),
      }),
      ('PUT', ['units', final id]) => _updateUnit(request, id),
      ('POST', ['units', final id, 'archive']) => _archive(id, true),
      ('POST', ['units', final id, 'restore']) => _archive(id, false),
      ('GET', ['units', final id, 'members']) => _json(200, [
        for (final userId in members[id] ?? const <String>{})
          {
            'userId': userId,
            'displayName': people[userId]!.$1,
            'email': people[userId]!.$2,
            'joinedAt': '2026-11-12T10:00:00Z',
          },
      ]),
      ('PUT', ['units', final id, 'members', final userId]) => _member(
        id,
        userId,
        true,
      ),
      ('DELETE', ['units', final id, 'members', final userId]) => _member(
        id,
        userId,
        false,
      ),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  List<Map<String, Object?>> _sortedRoles() => [...roles]
    ..sort((a, b) => (b['position']! as int).compareTo(a['position']! as int));

  Map<String, Object?> _body(http.Request request) =>
      jsonDecode(request.body) as Map<String, Object?>;

  http.Response _createRole(http.Request request) {
    final role = {
      ..._body(request),
      'id': 'r${_next++}',
      'icon': '',
      'position': 1,
      'default': false,
    };
    roles.add(role);
    return _json(201, role);
  }

  http.Response _updateRole(http.Request request, String id) {
    final role = roles.firstWhere((role) => role['id'] == id)
      ..addAll(_body(request));
    return _json(200, role);
  }

  http.Response _deleteRole(String id) {
    roles.removeWhere((role) => role['id'] == id);
    return http.Response('', 204);
  }

  http.Response _moveRole(http.Request request, String id) {
    final position = _body(request)['position']! as int;
    moves.add((id, position));
    final role = roles.firstWhere((role) => role['id'] == id);
    final from = role['position']! as int;
    for (final other in roles) {
      final at = other['position']! as int;
      if (other['default'] == false && other != role) {
        if (position > from && at > from && at <= position) {
          other['position'] = at - 1;
        } else if (position < from && at >= position && at < from) {
          other['position'] = at + 1;
        }
      }
    }
    role['position'] = position;
    return http.Response('', 204);
  }

  http.Response _assign(http.Request request) {
    final body = _body(request);
    final assignment = {
      'id': 'a${_next++}',
      'roleId': body['roleId'],
      'userId': body['userId'],
      'scope': 'tenant',
      'expiresAt': null,
      'grantedBy': 'u1',
      'grantedAt': '2026-11-12T10:00:00Z',
    };
    assignments.add(assignment);
    return _json(201, assignment);
  }

  http.Response _unassign(String id) {
    assignments.removeWhere((item) => item['id'] == id);
    return http.Response('', 204);
  }

  Map<String, Object?> _unit(String id) =>
      units.firstWhere((unit) => unit['id'] == id);

  List<Map<String, Object?>> _path(String id) {
    final path = <Map<String, Object?>>[];
    var parent = _unit(id)['parentId'] as String?;
    while (parent != null) {
      final unit = _unit(parent);
      path.insert(0, unit);
      parent = unit['parentId'] as String?;
    }
    return path;
  }

  http.Response _createUnit(http.Request request) {
    final body = _body(request);
    final unit = {
      'id': 'unit${_next++}',
      'parentId': body['parentId'],
      'kind': body['kind'],
      'name': body['name'],
      'code': body['code'],
      'archivedAt': null,
      'createdAt': '2026-11-12T10:00:00Z',
    };
    units.add(unit);
    return _json(201, unit);
  }

  http.Response _updateUnit(http.Request request, String id) {
    final unit = _unit(id)..addAll(_body(request));
    return _json(200, unit);
  }

  http.Response _archive(String id, bool archived) {
    _unit(id)['archivedAt'] = archived ? '2026-11-12T10:00:00Z' : null;
    return http.Response('', 204);
  }

  http.Response _member(String id, String userId, bool add) {
    final set = members.putIfAbsent(id, () => {});
    add ? set.add(userId) : set.remove(userId);
    return http.Response('', 204);
  }

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
