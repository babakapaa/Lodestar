import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeTeams {
  FakeTeams(this._role);

  final String Function() _role;
  final List<Map<String, Object?>> teams = [];
  final Map<String, List<String>> invited = {};
  final List<Map<String, Object?>> history = [];
  final Map<String, String> seeking = {};
  Map<String, Object?>? proposal;
  Map<String, Object?> card = {
    'interests': '',
    'skills': '',
    'hours': '',
    'role': '',
  };
  int _next = 1;

  static const me = 'u1';
  static const people = {
    'u1': {
      'userId': 'u1',
      'displayName': 'Айгерим Сапарова',
      'email': 'aigerim@kaznu.kz',
    },
    'u2': {
      'userId': 'u2',
      'displayName': 'Ерлан Ким',
      'email': 'erlan@kaznu.kz',
    },
    'u4': {
      'userId': 'u4',
      'displayName': 'Мадина Омарова',
      'email': 'madina@kaznu.kz',
    },
  };

  static const scope = {
    'kind': 'study.assignment',
    'id': 'a1',
    'title': 'Эссе о клетке',
    'context': 'Биология',
    'link': '/study/assignments/a1',
    'minSize': 2,
    'maxSize': 3,
    'selfFormed': true,
    'freeLeaveUntil': null,
    'dueAt': '2030-12-01T18:00:00Z',
    'archived': false,
  };

  bool get _leads => _role() != 'student';

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('GET', ['teams']) => _json(200, [
        for (final t in teams)
          if (_members(t).contains(me)) _card(t),
      ]),
      ('GET', ['invitations']) => _json(200, [
        for (final t in teams)
          if (invited[t['id']]?.contains(me) ?? false) _card(t),
      ]),
      ('GET', ['card']) => _json(200, card),
      ('PUT', ['card']) => _saveCard(request),
      ('GET', ['scopes', 'study.assignment', 'a1']) => _json(200, _view()),
      ('POST', ['scopes', 'study.assignment', 'a1', 'teams']) => _create(
        request,
      ),
      ('PUT', ['scopes', 'study.assignment', 'a1', 'seeking']) => _seek(
        request,
      ),
      ('POST', ['scopes', 'study.assignment', 'a1', 'distribution']) =>
        _distribute(),
      ('POST', ['scopes', 'study.assignment', 'a1', 'suggestions']) => _json(
        200,
        {
          'teams': [
            {
              'members': [people['u1'], people['u2']],
              'why': 'Оба любят микроскопию',
            },
          ],
          'without': <Object?>[],
        },
      ),
      ('GET', ['teams', final id]) => _details(id),
      ('PATCH', ['teams', final id]) => _change(id, request, (t, body) {
        t['name'] = body['name'];
      }),
      ('POST', ['teams', final id, 'invitations']) => _change(
        id,
        request,
        (t, body) =>
            invited.putIfAbsent(id, () => []).add(body['userId']! as String),
      ),
      ('PUT', ['teams', final id, 'members', final user]) => _change(
        id,
        request,
        (t, _) {
          for (final other in teams) {
            _members(other).remove(user);
          }
          _members(t).add(user);
        },
      ),
      ('DELETE', ['teams', final id, 'members', final user]) => _change(
        id,
        request,
        (t, _) => _members(t).remove(user),
      ),
      ('POST', ['teams', final id, 'leave']) => _leave(id),
      ('POST', ['teams', final id, 'proposals']) => _propose(id, request),
      ('POST', ['proposals', _, 'votes']) => _vote(request),
      _ => _json(404, {'code': 'team_not_found'}),
    };
  }

  List<String> _members(Map<String, Object?> team) =>
      team['members']! as List<String>;

  Map<String, Object?>? _team(String id) {
    for (final t in teams) {
      if (t['id'] == id) {
        return t;
      }
    }
    return null;
  }

  Map<String, Object?> _teamJson(Map<String, Object?> t) => {
    'id': t['id'],
    'scopeKind': 'study.assignment',
    'scopeId': 'a1',
    'name': t['name'],
    'archivedAt': null,
    'members': [
      for (final m in _members(t))
        {
          ...people[m]!,
          'label': m == me ? (t['label'] ?? '') : '',
          'joinedAt': '2026-11-12T10:00:00Z',
        },
    ],
  };

  Map<String, Object?> _card(Map<String, Object?> t) => {
    'team': _teamJson(t),
    'scope': scope,
    'invitedBy': null,
  };

  Map<String, Object?> _view() {
    final placed = {for (final t in teams) ..._members(t)};
    final mine = teams.where((t) => _members(t).contains(me)).toList();
    return {
      'scope': scope,
      'role': _leads ? 'lead' : 'student',
      'myTeam': mine.isEmpty || _leads ? null : _teamJson(mine.first),
      'solo': false,
      'seeking': seeking.containsKey(me),
      'seekers': [
        for (final e in seeking.entries) {...people[e.key]!, 'note': e.value},
      ],
      'teams': _leads ? [for (final t in teams) _teamJson(t)] : <Object?>[],
      'unplaced': [
        for (final p in people.keys)
          if (!placed.contains(p) && (_leads || seeking.containsKey(p)))
            people[p],
      ],
      'solos': <Object?>[],
      'requests': <Object?>[],
    };
  }

  http.Response _details(String id) {
    final t = _team(id);
    if (t == null) {
      return _json(404, {'code': 'team_not_found'});
    }
    return _json(200, {
      ..._teamJson(t),
      'scope': scope,
      'manage': _leads,
      'inviteCode': 'JOIN-$id',
      'invited': [for (final u in invited[id] ?? const <String>[]) people[u]],
      'history': [
        for (final h in history)
          if (h['team'] == id) h,
      ],
      'request': null,
      'canLeaveFreely': true,
      'proposal': proposal?['teamId'] == id ? proposal : null,
    });
  }

  http.Response _create(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final id = 't${_next++}';
    final members = _leads
        ? [for (final m in body['members']! as List<Object?>) m! as String]
        : [me];
    teams.add({'id': id, 'name': body['name'], 'members': members});
    seeking.remove(me);
    history.add({
      'team': id,
      'person': null,
      'event': 'created',
      'actor': people[me],
      'at': '2026-11-13T10:00:00Z',
    });
    return _details(id);
  }

  http.Response _seek(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    if (body['seeking'] == true) {
      seeking[me] = body['note']! as String;
    } else {
      seeking.remove(me);
    }
    return _json(200, _view());
  }

  http.Response _distribute() {
    final placed = {for (final t in teams) ..._members(t)};
    final rest = [
      for (final p in people.keys)
        if (!placed.contains(p)) p,
    ];
    teams.add({'id': 't${_next++}', 'name': 'Айгерим, Ерлан', 'members': rest});
    return _json(200, _view());
  }

  http.Response _change(
    String id,
    http.Request request,
    void Function(Map<String, Object?> team, Map<String, Object?> body) apply,
  ) {
    final t = _team(id);
    if (t == null) {
      return _json(404, {'code': 'team_not_found'});
    }
    apply(
      t,
      request.body.isEmpty
          ? const {}
          : jsonDecode(request.body) as Map<String, Object?>,
    );
    return _details(id);
  }

  http.Response _leave(String id) {
    _members(_team(id)!).remove(me);
    return _json(200, {'left': true, 'request': null});
  }

  http.Response _propose(String id, http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final t = _team(id)!;
    proposal = {
      'id': 'p${_next++}',
      'teamId': id,
      'proposedBy': people[me],
      'work': body['work'],
      'note': body['note'],
      'status': 'open',
      'silentAt': '2030-11-14T10:00:00Z',
      'votes': [
        for (final m in _members(t))
          {...people[m]!, 'decision': m == me ? 'agreed' : 'waiting'},
      ],
    };
    return _details(id);
  }

  http.Response _vote(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final p = proposal!;
    p['votes'] = [
      for (final v in p['votes']! as List<Object?>)
        if ((v! as Map<String, Object?>)['userId'] == me)
          {
            ...v as Map<String, Object?>,
            'decision': body['agree'] == true ? 'agreed' : 'objected',
            'comment': body['comment'],
          }
        else
          v,
    ];
    return _details(p['teamId']! as String);
  }

  http.Response _saveCard(http.Request request) {
    card = jsonDecode(request.body) as Map<String, Object?>;
    return _json(200, card);
  }

  http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
