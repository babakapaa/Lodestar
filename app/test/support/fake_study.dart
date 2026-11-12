import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeStudy {
  String role = 'student';
  String? submitError;
  final List<Map<String, Object?>> assignments = [
    {
      'id': 'a1',
      'courseId': 'c1',
      'title': 'Эссе о клетке',
      'description': 'Опишите строение клетки.',
      'dueAt': '2030-12-01T18:00:00Z',
      'grading': 'points',
      'maxPoints': 10,
      'publishedAt': '2026-11-10T10:00:00Z',
    },
  ];
  final List<Map<String, Object?>> submissions = [];
  final List<Map<String, Object?>> lessons = [];
  final List<Map<String, Object?>> deadlines = [];
  final Map<String, Map<String, Object?>> grades = {};
  int _next = 1;

  static const _student = {
    'userId': 'u2',
    'displayName': 'Ерлан Ким',
    'email': 'erlan@kaznu.kz',
  };

  Map<String, Object?> get _course => {
    'id': 'c1',
    'title': 'Биология',
    'code': 'BIO-101',
    'description': 'Основы клеточной биологии',
    'archivedAt': null,
    'role': role,
  };

  bool get _teaches => role != 'student';

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('GET', ['courses']) => _json(200, [_course]),
      ('GET', ['courses', 'c1']) => _json(200, {
        ..._course,
        'canManage': false,
        'teachers': [
          {
            'userId': 'u3',
            'displayName': 'Дана Ахметова',
            'email': 'dana@kaznu.kz',
          },
        ],
        'units': [
          {'id': 'g1', 'name': 'БИО-21'},
        ],
      }),
      ('GET', ['courses', 'c1', 'assignments']) => _json(200, [
        for (final a in assignments)
          if (_teaches || a['publishedAt'] != null) _item(a),
      ]),
      ('POST', ['courses', 'c1', 'assignments']) => _create(request),
      ('GET', ['assignments', final id]) => _details(id),
      ('PUT', ['assignments', final id]) => _update(request, id),
      ('POST', ['assignments', final id, 'submissions']) => _submit(
        request,
        id,
      ),
      ('GET', ['assignments', final id, 'work']) => _json(200, [_work(id)]),
      ('GET', ['assignments', final id, 'students', 'u2', 'submissions']) =>
        _json(200, _mine(id)),
      ('PUT', ['assignments', final id, 'grades', 'u2']) => _grade(request, id),
      ('GET', ['agenda']) => _json(200, {
        'date': '2026-11-12',
        'timeZone': 'Asia/Almaty',
        'lessons': lessons,
        'deadlines': deadlines,
      }),
      ('GET', ['schedule']) => _json(200, lessons),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  Map<String, Object?> _assignment(String id) =>
      assignments.firstWhere((a) => a['id'] == id);

  List<Map<String, Object?>> _mine(String id) => [
    for (final s in submissions.reversed)
      if (s['assignmentId'] == id) s,
  ];

  Map<String, Object?> _item(Map<String, Object?> a) {
    final mine = _mine(a['id']! as String);
    return {
      ...a,
      'mine': _teaches
          ? null
          : {
              'submittedAt': mine.isEmpty ? null : mine.first['submittedAt'],
              'grade': null,
            },
      'counts': _teaches
          ? {'students': 1, 'submitted': mine.isEmpty ? 0 : 1, 'graded': 0}
          : null,
    };
  }

  Map<String, Object?> _work(String id) {
    final mine = _mine(id);
    return {
      'student': _student,
      'latest': mine.isEmpty ? null : mine.first,
      'attempts': mine.length,
      'grade': grades[id],
    };
  }

  http.Response _details(String id) {
    final a = _assignment(id);
    return _json(200, {
      ...a,
      'courseTitle': 'Биология',
      'role': role,
      'files': <Object?>[],
      'mine': _teaches ? null : _work(id),
      'counts': _teaches
          ? {'students': 1, 'submitted': _mine(id).isEmpty ? 0 : 1, 'graded': 0}
          : null,
    });
  }

  Map<String, Object?> _fromInput(Map<String, Object?> body) => {
    'title': body['title'],
    'description': body['description'],
    'dueAt': body['dueAt'],
    'grading': body['grading'],
    'maxPoints': body['maxPoints'],
    'publishedAt': body['published'] == true ? '2026-11-12T10:00:00Z' : null,
  };

  http.Response _create(http.Request request) {
    final a = {
      'id': 'a${++_next}',
      'courseId': 'c1',
      ..._fromInput(_body(request)),
    };
    assignments.add(a);
    return _json(201, a);
  }

  http.Response _update(http.Request request, String id) {
    final a = _assignment(id)..addAll(_fromInput(_body(request)));
    return _json(200, a);
  }

  http.Response _submit(http.Request request, String id) {
    final error = submitError;
    if (error != null) {
      return _json(409, {'code': error});
    }
    final body = _body(request);
    for (final s in submissions) {
      if (s['clientId'] == body['clientId']) {
        return _json(200, s);
      }
    }
    final submission = {
      'id': 's${++_next}',
      'assignmentId': id,
      'clientId': body['clientId'],
      'body': body['body'],
      'submittedAt': '2026-11-12T10:00:00Z',
      'files': <Object?>[],
    };
    submissions.add(submission);
    return _json(201, submission);
  }

  http.Response _grade(http.Request request, String id) {
    final grade = {
      ..._body(request),
      'assignmentId': id,
      'studentId': 'u2',
      'gradedAt': '2026-11-12T11:00:00Z',
    };
    grades[id] = grade;
    return _json(200, grade);
  }

  Map<String, Object?> _body(http.Request request) =>
      jsonDecode(request.body) as Map<String, Object?>;

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
