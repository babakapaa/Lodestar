import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeSandbox {
  final Map<String, Object?> term = {
    'id': 's1',
    'name': 'Осень',
    'startsOn': '2026-09-01',
    'endsOn': '2026-12-31',
  };
  final List<Map<String, Object?>> runs = [];
  final List<Map<String, Object?>> started = [];
  final List<String> stopped = [];
  final List<Map<String, Object?>> transfers = [];
  final List<Map<String, Object?>> decisions = [];
  final List<Map<String, Object?>> topups = [];
  final List<Map<String, Object?>> nodes = [];
  final List<Map<String, Object?>> templates = [
    {
      'id': 'tp1',
      'name': 'PyTorch',
      'about': '',
      'image': 'pytorch/pytorch:2.5.1',
      'archived': false,
      'starter': false,
      'createdAt': '2026-09-01T00:00:00Z',
    },
  ];
  final List<Map<String, Object?>> classes = [
    {
      'id': 'c-t4',
      'name': 'T4',
      'memoryGb': 16,
      'rate': 2,
      'archived': false,
      'createdAt': '2026-09-01T00:00:00Z',
    },
  ];
  int _next = 1;

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('GET', ['quota']) => _json(200, _quota()),
      ('POST', ['transfers']) => _record(transfers, request),
      ('GET', ['topups']) => _json(200, topups),
      ('POST', ['topups']) => _askTopup(request),
      ('POST', ['topups', final id, 'decision']) => _decideTopup(request, id),
      ('POST', ['estimate']) => _json(200, _estimate(request)),
      ('GET', ['runs']) => _json(200, runs),
      ('GET', ['runs', 'held']) => _json(200, [
        for (final r in runs)
          if (r['status'] == 'held') r,
      ]),
      ('GET', ['runs', 'export']) => http.Response.bytes([1, 2, 3], 200),
      ('POST', ['runs']) => _start(request),
      ('GET', ['runs', final id]) => _run(id),
      ('POST', ['runs', final id, 'stop']) => _stop(id),
      ('POST', ['runs', final id, 'rerun']) => _rerun(id),
      ('POST', ['runs', final id, 'hold-decision']) => _hold(request, id),
      ('GET', ['templates']) => _json(200, templates),
      ('GET', ['classes']) => _json(200, classes),
      ('GET', ['terms']) => _json(200, [term]),
      ('GET', ['nodes']) => _json(200, nodes),
      ('POST', ['nodes']) => _node(request),
      ('DELETE', ['nodes', final id]) => _deleteNode(id),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  Map<String, Object?> run(String id, String kind, String status) => {
    'id': id,
    'userId': 'u1',
    'displayName': 'Айгерим Сапарова',
    'teamId': null,
    'teamName': '',
    'kind': kind,
    'title': 'Запуск $id',
    'command': kind == 'jupyter' ? '' : 'python app.py',
    'templateId': 'tp1',
    'templateName': 'PyTorch',
    'classId': 'c-t4',
    'className': 'T4',
    'rate': 2,
    'maxHours': 4,
    'workKind': 'study.assignment',
    'workId': 'a1',
    'workTitle': 'Нейросеть для МРТ',
    'status': status,
    'position': status == 'queued' ? 1 : null,
    'url': status == 'running' ? 'http://gpu-01:49153/' : '',
    'exitCode': null,
    'reason': '',
    'charged': null,
    'queuedAt': '2026-11-13T08:00:00Z',
    'startedAt': null,
    'endedAt': null,
    'heldAt': null,
    'holdNote': '',
    'refunded': false,
    'gpuLoad': null,
    'log': '',
  };

  Map<String, Object?> _quota() => {
    'term': term,
    'lines': [
      {'source': 'sandbox.base', 'title': '', 'hours': 20},
      {'source': 'league.division', 'title': 'Вершина', 'hours': 10},
    ],
    'hours': 30,
    'reserved': 4,
    'spent': 0,
    'forecast': 30,
    'teams': [
      {
        'id': 'team-1',
        'name': 'Нейроны',
        'hours': 2,
        'reserved': 0,
        'lines': <Object?>[],
      },
    ],
  };

  Map<String, Object?> _estimate(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final hours = (body['maxHours'] as num?)?.toDouble() ?? 4;
    return {
      'maxHours': hours,
      'starter': false,
      'available': 26,
      'night': {'from': 23, 'to': 7, 'discount': 50, 'timeZone': 'Asia/Almaty'},
      'options': [
        {
          'classId': 'c-t4',
          'className': 'T4',
          'memoryGb': 16,
          'rate': 2,
          'cost': hours * 2,
          'nightCost': hours,
          'freeGpus': 1,
          'queued': 0,
          'enough': true,
        },
        {
          'classId': 'c-a100',
          'className': 'A100',
          'memoryGb': 80,
          'rate': 8,
          'cost': hours * 8,
          'nightCost': hours * 4,
          'freeGpus': 0,
          'queued': 0,
          'enough': hours * 8 <= 26,
        },
      ],
      'suggested': 'c-t4',
    };
  }

  http.Response _start(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    started.add(body);
    final created = run('r${_next++}', body['kind']! as String, 'queued')
      ..['title'] = body['title']
      ..['classId'] = body['classId'];
    runs.insert(0, created);
    return _json(201, created);
  }

  http.Response _run(String id) {
    final found = runs.where((r) => r['id'] == id).firstOrNull;
    return found == null
        ? _json(404, {'code': 'run_not_found'})
        : _json(200, found);
  }

  http.Response _stop(String id) {
    stopped.add(id);
    final found = runs.firstWhere((r) => r['id'] == id);
    found['status'] = 'stopped';
    return _json(200, found);
  }

  http.Response _rerun(String id) {
    final old = runs.firstWhere((r) => r['id'] == id);
    final created = run('r${_next++}', old['kind']! as String, 'queued')
      ..['title'] = old['title'];
    runs.insert(0, created);
    return _json(201, created);
  }

  http.Response _hold(http.Request request, String id) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    decisions.add({'run': id, ...body});
    final found = runs.firstWhere((r) => r['id'] == id);
    found['status'] = body['release'] == true ? 'running' : 'stopping';
    return _json(200, found);
  }

  http.Response _askTopup(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    topups.add({
      'id': 'tu${_next++}',
      'userId': 'u1',
      'displayName': 'Айгерим Сапарова',
      'title': 'Нейросеть для МРТ',
      'status': 'pending',
      'mine': true,
      'createdAt': '2026-11-13T08:00:00Z',
      ...body,
    });
    return _json(201, topups.last);
  }

  http.Response _decideTopup(http.Request request, String id) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final found = topups.firstWhere((t) => t['id'] == id);
    found['status'] = body['approve'] == true ? 'approved' : 'declined';
    found['granted'] = body['hours'];
    decisions.add({'topup': id, ...body});
    return _json(200, found);
  }

  http.Response _node(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    final created = {
      'id': 'n${_next++}',
      'name': body['name'],
      'classId': body['classId'],
      'className': 'T4',
      'gpus': body['gpus'],
      'busy': 0,
      'seenAt': null,
      'online': false,
      'createdAt': '2026-11-13T08:00:00Z',
    };
    nodes.add(created);
    return _json(201, {...created, 'token': 'lsn_t1_secret'});
  }

  http.Response _deleteNode(String id) {
    nodes.removeWhere((n) => n['id'] == id);
    return http.Response('', 204);
  }

  http.Response _record(List<Map<String, Object?>> into, http.Request request) {
    into.add(jsonDecode(request.body) as Map<String, Object?>);
    return http.Response('', 204);
  }

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
