import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeAttendance {
  final List<Map<String, Object?>> checks = [];
  final List<Map<String, Object?>> marked = [];
  final Map<String, Map<String, Object?>> sessions = {};
  final List<Map<String, Object?>> open = [];
  int _next = 1;

  final List<Map<String, Object?>> lessons = [
    {
      'seriesId': 'ser1',
      'date': '2026-09-01',
      'startsAt': '2026-09-01T09:00:00Z',
      'room': '301',
      'cancelled': false,
      'sessionId': 'a0',
      'open': false,
      'counts': _counts(present: 1, absent: 1),
    },
    {
      'seriesId': 'ser1',
      'date': '2026-09-08',
      'startsAt': '2026-09-08T09:00:00Z',
      'room': '301',
      'cancelled': false,
      'sessionId': null,
      'open': false,
      'counts': _counts(),
    },
  ];

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('POST', ['lessons', final series, final date, 'open']) => _open(
        series,
        date,
      ),
      ('GET', ['sessions', final id]) => _session(id),
      ('GET', ['sessions', _, 'qr']) => _json(200, {
        'token': 'tok',
        'expiresAt': '2999-01-01T00:00:00Z',
      }),
      ('POST', ['sessions', final id, 'code']) => _session(id),
      ('PUT', ['sessions', final id, 'marks', final user]) => _mark(
        request,
        id,
        user,
      ),
      ('POST', ['sessions', final id, 'close']) => _close(id),
      ('POST', ['check']) => _check(request),
      ('GET', ['open']) => _json(200, open),
      ('GET', ['courses', final id]) => _json(200, _journal(id)),
      ('GET', ['courses', _, 'export']) => http.Response.bytes([1, 2, 3], 200),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  static Map<String, Object?> _counts({
    int present = 0,
    int absent = 0,
    int unmarked = 0,
  }) => {
    'present': present,
    'late': 0,
    'absent': absent,
    'excused': 0,
    'unmarked': unmarked,
  };

  Map<String, Object?> _journal(String id) => {
    'courseId': id,
    'courseTitle': 'Биология',
    'from': '2026-05-01',
    'to': '2026-09-08',
    'absentAlert': 3,
    'lessons': lessons,
    'students': [
      {
        'userId': 'u1',
        'name': 'Айгерим Сапарова',
        'counts': _counts(present: 1),
        'rate': 1.0,
        'streak': 0,
        'longAbsent': false,
      },
      {
        'userId': 'u2',
        'name': 'Ерлан Ахметов',
        'counts': _counts(absent: 3),
        'rate': 0.0,
        'streak': 3,
        'longAbsent': true,
      },
    ],
  };

  http.Response _open(String series, String date) {
    final id = 'a${_next++}';
    sessions[id] = {
      'id': id,
      'seriesId': series,
      'date': date,
      'courseId': 'c1',
      'courseTitle': 'Биология',
      'room': '301',
      'startsAt': '${date}T09:00:00Z',
      'endsAt': '${date}T10:20:00Z',
      'open': true,
      'code': '4821',
      'codeUntil': '${date}T09:05:00Z',
      'methods': {'qr': true, 'code': true},
      'marks': [
        {'userId': 'u1', 'name': 'Айгерим Сапарова', 'status': null},
        {'userId': 'u2', 'name': 'Ерлан Ахметов', 'status': null},
      ],
    };
    return _session(id);
  }

  http.Response _session(String id) {
    final s = sessions[id];
    if (s == null) {
      return _json(404, {'code': 'attendance_session_not_found'});
    }
    final marks = s['marks']! as List<Map<String, Object?>>;
    final present = marks.where((m) => m['status'] == 'present').length;
    final absent = marks.where((m) => m['status'] == 'absent').length;
    final unmarked = marks.where((m) => m['status'] == null).length;
    return _json(200, {
      ...s,
      'counts': _counts(present: present, absent: absent, unmarked: unmarked),
    });
  }

  http.Response _mark(http.Request request, String id, String user) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    marked.add({'session': id, 'user': user, ...body});
    for (final m in sessions[id]!['marks']! as List<Map<String, Object?>>) {
      if (m['userId'] == user) {
        m['status'] = body['status'];
      }
    }
    return _session(id);
  }

  http.Response _close(String id) {
    final s = sessions[id]!;
    s['open'] = false;
    s.remove('code');
    for (final m in s['marks']! as List<Map<String, Object?>>) {
      m['status'] ??= 'absent';
    }
    return _session(id);
  }

  http.Response _check(http.Request request) {
    final body = jsonDecode(request.body) as Map<String, Object?>;
    checks.add(body);
    if (body['code'] != null && body['code'] != '4821') {
      return _json(400, {'code': 'wrong_attendance_code'});
    }
    for (final o in open) {
      o['status'] = 'present';
    }
    return _json(200, {
      'sessionId': 'a1',
      'courseTitle': 'Биология',
      'startsAt': '2026-09-08T09:00:00Z',
      'status': 'present',
    });
  }

  http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
