import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeLeague {
  FakeLeague(this._role);

  final String Function() _role;
  Map<String, Object?> standing = {
    'points': 120,
    'division': 0,
    'toNext': 30,
    'growth': 20,
    'display': 'name',
    'nickname': '',
    'gpuHours': 0,
  };
  List<Map<String, Object?>> divisions = [
    {'name': 'Бронза', 'from': 0, 'gpuHours': 0},
    {'name': 'Серебро', 'from': 150, 'gpuHours': 10},
  ];
  bool confirm = true;
  int weight = 100;
  final List<Map<String, Object?>> disputes = [];
  final List<Map<String, Object?>> pending = [
    {
      'id': 'm1',
      'seasonId': 's1',
      'userId': 'u2',
      'from': 0,
      'to': 1,
      'points': 160,
      'status': 'pending',
      'createdAt': '2026-11-12T10:00:00Z',
      'person': erlan,
    },
  ];
  final List<Map<String, Object?>> adjustments = [];
  final List<Map<String, Object?>> tasks = [];
  final List<Map<String, Object?>> awards = [];
  final List<Map<String, Object?>> grants = [];
  Map<String, Object?>? review;
  int _next = 1;

  static const erlan = {
    'userId': 'u2',
    'displayName': 'Ерлан Ким',
    'email': 'erlan@kaznu.kz',
  };

  static const entry = {
    'id': 'e1',
    'seasonId': 's1',
    'userId': 'u1',
    'kind': 'work',
    'workKind': 'study.assignment',
    'workId': 'a1',
    'title': 'Эссе о клетке',
    'score': 8,
    'max': 10,
    'weight': 100,
    'points': 80,
    'reason': '',
    'createdAt': '2026-11-11T10:00:00Z',
    'updatedAt': '2026-11-11T10:00:00Z',
  };

  bool get _leads => _role() != 'student';

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('GET', ['arenas', 'study.course', 'c1']) => _json(200, _view()),
      ('PUT', ['leagues', 'l1', 'display']) => _display(request),
      ('PUT', ['leagues', 'l1', 'rules']) => _rules(request),
      ('POST', ['leagues', 'l1', 'disputes']) => _dispute(request),
      ('POST', ['leagues', 'l1', 'adjustments']) => _adjust(request),
      ('POST', ['leagues', 'l1', 'tasks']) => _createTask(request),
      ('POST', ['leagues', 'l1', 'grants']) => _grant(request),
      ('GET', ['leagues', 'l1', 'players', 'u2']) => _json(200, _player()),
      ('POST', ['transitions', 'm1']) => _decide(),
      ('GET', ['tasks', final id]) => _task(id),
      ('POST', ['solutions', 'sol1', 'review']) => _review(request),
      ('GET', ['awards']) => _json(200, awards),
      ('POST', ['awards']) => _createAward(request),
      _ => _json(404, {'code': 'league_not_found'}),
    };
  }

  Map<String, Object?> _view() => {
    'id': 'l1',
    'kind': 'study.course',
    'originId': 'c1',
    'title': 'Биология',
    'role': _leads ? 'lead' : 'student',
    'season': {
      'id': 's1',
      'number': 2,
      'startedAt': '2026-09-01T00:00:00Z',
      'endedAt': null,
    },
    'divisions': divisions,
    'confirm': confirm,
    'formula': [
      {
        'kind': 'study.assignment',
        'id': 'a1',
        'title': 'Эссе о клетке',
        'project': false,
        'weight': weight,
      },
    ],
    'me': _leads ? null : standing,
    'table': _leads
        ? <Object?>[]
        : [
            {
              'place': 1,
              'name': 'Ерлан Ким',
              'display': 'name',
              'points': 140,
              'growth': 50,
              'me': false,
              'topGrowth': true,
            },
            {
              'place': 2,
              'name': 'Айгерим Сапарова',
              'display': standing['display'],
              'points': standing['points'],
              'growth': 20,
              'me': true,
              'topGrowth': false,
            },
          ],
    'entries': _leads ? <Object?>[] : [entry],
    'transitions': <Object?>[],
    'disputes': _leads ? <Object?>[] : disputes,
    'players': _leads
        ? [
            {
              'person': erlan,
              'division': 0,
              'points':
                  160 +
                  adjustments.fold<int>(
                    0,
                    (sum, a) => sum + (a['points']! as int),
                  ),
              'growth': 50,
              'topGrowth': true,
              'display': 'name',
              'nickname': '',
            },
          ]
        : <Object?>[],
    'pending': _leads ? pending : <Object?>[],
    'openDisputes': _leads ? disputes : <Object?>[],
    'tasks': tasks,
    'grants': _leads ? grants : <Object?>[],
  };

  Map<String, Object?> _player() => {
    'person': erlan,
    'standing': {...standing, 'points': 160},
    'entries': [
      for (final a in adjustments)
        {
          ...entry,
          'id': 'adj${a['points']}',
          'kind': 'adjustment',
          'title': a['reason'],
          'points': a['points'],
          'reason': a['reason'],
        },
    ],
    'transitions': <Object?>[],
    'disputes': <Object?>[],
    'grants': [
      for (final g in grants)
        if (g['userId'] == 'u2') g,
    ],
  };

  http.Response _display(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    standing = {
      ...standing,
      'display': body['display'],
      'nickname': body['nickname'],
    };
    return _json(200, _view());
  }

  http.Response _rules(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    divisions = [
      for (final d in body['divisions']! as List<Object?>)
        d! as Map<String, Object?>,
    ];
    confirm = body['confirm']! as bool;
    weight =
        ((body['weights']! as List<Object?>).single!
                as Map<String, Object?>)['weight']!
            as int;
    return _json(200, _view());
  }

  http.Response _dispute(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final dispute = {
      'id': 'd${_next++}',
      'userId': 'u1',
      'entryId': body['entryId'],
      'transitionId': body['transitionId'],
      'text': body['text'],
      'status': 'open',
      'answer': '',
      'createdAt': '2026-11-13T10:00:00Z',
    };
    disputes.add(dispute);
    return _json(201, dispute);
  }

  http.Response _adjust(http.Request request) {
    adjustments.add(jsonDecode(request.body) as Map<String, Object?>);
    return _json(200, _player());
  }

  http.Response _decide() {
    pending.clear();
    return _json(200, {'id': 'm1', 'status': 'applied'});
  }

  Map<String, Object?> _taskJson(Map<String, Object?> body, String id) => {
    'id': id,
    'leagueId': 'l1',
    'title': body['title'],
    'about': body['about'],
    'partner': body['partner'],
    'points': body['points'],
    'minSize': body['minSize'],
    'maxSize': body['maxSize'],
    'selfFormed': body['selfFormed'],
    'dueAt': body['dueAt'],
    'closedAt': null,
    'createdAt': '2026-11-13T10:00:00Z',
  };

  http.Response _createTask(http.Request request) {
    final task = _taskJson(
      jsonDecode(request.body) as Map<String, Object?>,
      'task${_next++}',
    );
    tasks.add(task);
    return _json(201, task);
  }

  http.Response _task(String id) {
    final task = tasks.where((t) => t['id'] == id).firstOrNull;
    if (task == null) {
      return _json(404, {'code': 'league_task_not_found'});
    }
    final reviewed = review;
    return _json(200, {
      'task': task,
      'league': 'Биология',
      'kind': 'study.course',
      'originId': 'c1',
      'role': _leads ? 'lead' : 'student',
      'people': [
        erlan,
        {
          'userId': 'u4',
          'displayName': 'Мадина Омарова',
          'email': 'madina@kaznu.kz',
        },
      ],
      'solutions': _leads
          ? [
              {
                'id': 'sol1',
                'taskId': id,
                'teamId': 't9',
                'teamName': 'Микробы',
                'members': ['u2', 'u4'],
                'work': {
                  'text': 'Бот для поиска статей',
                  'links': ['https://example.org/bot'],
                },
                'votes': <Object?>[],
                'similarTo': 'old',
                'similarity': 0.93,
                'explanation': '',
                'points': reviewed?['points'],
                'passed': reviewed?['passed'] ?? <Object?>[],
                'comment': reviewed?['comment'] ?? '',
                'reviewedAt': reviewed == null ? null : '2026-11-13T12:00:00Z',
                'createdAt': '2026-11-13T10:00:00Z',
                'similar': {
                  'id': 'old',
                  'taskTitle': 'Бот прошлого года',
                  'teamName': 'Старички',
                  'text': 'Бот для поиска литературы',
                },
              },
            ]
          : <Object?>[],
    });
  }

  http.Response _review(http.Request request) {
    review = jsonDecode(request.body) as Map<String, Object?>;
    return _json(200, {
      'id': 'sol1',
      'teamName': 'Микробы',
      'members': ['u2', 'u4'],
      'work': {'text': '', 'links': <Object?>[]},
      'reviewedAt': '2026-11-13T12:00:00Z',
      'createdAt': '2026-11-13T10:00:00Z',
    });
  }

  http.Response _createAward(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final award = {
      'id': 'aw${_next++}',
      'title': body['title'],
      'about': body['about'],
      'gpuHours': body['gpuHours'],
      'quantity': body['quantity'],
      'granted': 0,
      'archivedAt': null,
      'createdAt': '2026-11-13T10:00:00Z',
    };
    awards.add(award);
    return _json(201, award);
  }

  http.Response _grant(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final award = awards.firstWhere((a) => a['id'] == body['awardId']);
    final grant = {
      'id': 'g${_next++}',
      'awardId': award['id'],
      'title': award['title'],
      'gpuHours': award['gpuHours'],
      'userId': body['userId'],
      'seasonId': 's1',
      'reason': body['reason'],
      'grantedAt': '2026-11-13T10:00:00Z',
    };
    grants.add(grant);
    return _json(201, grant);
  }

  http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
