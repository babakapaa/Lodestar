import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeDefense {
  FakeDefense(this._grades);

  final Map<String, Map<String, Object?>> _grades;
  final List<Map<String, Object?>> slots = [];
  String? slotId;
  bool selected = true;
  String status = 'planned';
  List<Object?> questions = [];
  String summary = '';
  String? heldAt;
  int _next = 1;

  static const student = {
    'userId': 'u2',
    'displayName': 'Ерлан Ким',
    'email': 'erlan@kaznu.kz',
  };

  static const drafts = [
    {'kind': 'understanding', 'text': 'Что такое мембрана клетки?'},
    {'kind': 'path', 'text': 'Как вы нашли главный источник?'},
  ];

  http.Response? handle(
    http.Request request,
    List<String> segments,
    Map<String, Object?> assignment,
    bool teaches,
  ) {
    final id = assignment['id']! as String;
    return switch ((request.method, segments.skip(2).toList())) {
      ('GET', ['defense-slots']) => _json(200, [
        for (final s in slots) _slot(s, teaches),
      ]),
      ('POST', ['defense-slots']) => _addSlots(request, id),
      ('GET', ['defense']) => _json(200, _mine(assignment)),
      ('PUT', ['defense', 'booking']) => _book(request, assignment),
      ('DELETE', ['defense', 'booking']) => _cancel(assignment),
      ('GET', ['defenses']) => _json(200, [_row(id)]),
      ('POST', ['defenses', 'sample']) => _json(200, {'chosen': <Object?>[]}),
      ('GET', ['defenses', 'u2']) => _json(200, _defense(id)),
      ('PUT', ['defenses', 'u2']) => _save(request, id),
      ('PUT', ['defenses', 'u2', 'selection']) => _select(request, id),
      ('POST', ['defenses', 'u2', 'finish']) => _finish(request, id),
      ('POST', ['defenses', 'u2', 'questions-draft']) => _json(200, {
        'questions': drafts,
      }),
      ('POST', ['defenses', 'u2', 'protocol-draft']) => _json(200, {
        'summary': 'Студент уверенно объяснил строение клетки.',
      }),
      _ => null,
    };
  }

  http.Response? deleteSlot(String id) {
    final removed = slots.where((s) => s['id'] == id).toList();
    if (removed.isEmpty) {
      return null;
    }
    slots.removeWhere((s) => s['id'] == id);
    if (slotId == id) {
      slotId = null;
    }
    return http.Response('', 204);
  }

  Map<String, Object?> _slot(Map<String, Object?> s, bool teaches) {
    final taken = slotId == s['id'];
    return {
      ...s,
      'taken': taken ? 1 : 0,
      'mine': !teaches && taken,
      if (teaches) 'students': [if (taken) student],
    };
  }

  Map<String, Object?>? get _booked {
    for (final s in slots) {
      if (s['id'] == slotId) {
        return _slot(s, false);
      }
    }
    return null;
  }

  Map<String, Object?> _defense(String id) => {
    'assignmentId': id,
    'studentId': 'u2',
    'selected': selected,
    'status': status,
    'questions': questions,
    'summary': summary,
    'heldAt': heldAt,
  };

  Map<String, Object?> _mine(Map<String, Object?> assignment) {
    final format = assignment['defense'] as String?;
    return {
      'format': format,
      'required': format != 'sample' || selected,
      'recorded': false,
      'status': status,
      'slot': _booked,
      'protocol': status == 'planned'
          ? null
          : _defense(assignment['id']! as String),
    };
  }

  Map<String, Object?> _row(String id) => {
    'student': student,
    'selected': selected,
    'status': status,
    'slot': _booked,
    'grade': _grades[id],
  };

  http.Response _addSlots(http.Request request, String id) {
    final body = _body(request);
    final start = DateTime.parse(body['startsAt']! as String);
    final minutes = body['minutes']! as int;
    final added = [
      for (var i = 0; i < (body['count'] as int? ?? 1); i++)
        {
          'id': 'slot${_next++}',
          'assignmentId': id,
          'startsAt': start
              .add(Duration(minutes: minutes * i))
              .toUtc()
              .toIso8601String(),
          'minutes': minutes,
          'place': body['place'] ?? '',
          'capacity': body['capacity'] ?? 1,
        },
    ];
    slots.addAll(added);
    return _json(201, added);
  }

  http.Response _book(http.Request request, Map<String, Object?> assignment) {
    slotId = _body(request)['slotId'] as String?;
    return _json(200, _mine(assignment));
  }

  http.Response _cancel(Map<String, Object?> assignment) {
    slotId = null;
    return _json(200, _mine(assignment));
  }

  http.Response _save(http.Request request, String id) {
    final body = _body(request);
    questions = body['questions'] as List<Object?>? ?? [];
    summary = body['summary'] as String? ?? '';
    return _json(200, _defense(id));
  }

  http.Response _select(http.Request request, String id) {
    selected = _body(request)['selected'] as bool? ?? false;
    return _json(200, _defense(id));
  }

  http.Response _finish(http.Request request, String id) {
    final body = _body(request);
    status = body['status']! as String;
    heldAt = '2026-11-13T10:00:00Z';
    final grade = body['grade'];
    if (grade is Map<String, Object?>) {
      _grades[id] = {
        ...grade,
        'assignmentId': id,
        'studentId': 'u2',
        'gradedAt': '2026-11-13T10:00:00Z',
      };
    }
    return _json(200, _defense(id));
  }

  Map<String, Object?> _body(http.Request request) =>
      jsonDecode(request.body) as Map<String, Object?>;

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
