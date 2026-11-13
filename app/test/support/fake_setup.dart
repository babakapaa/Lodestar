import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeSetup {
  bool launched = true;
  bool manages = true;
  bool signed = false;
  final Map<String, String> statuses = {};
  final Map<String, Map<String, Object?>> answers = {};
  final Map<String, String?> assignees = {};
  final List<String> calls = [];

  static const _sections = [
    ('deploy', true),
    ('university', true),
    ('people', true),
    ('admins', true),
    ('privacy', false),
    ('legal', true),
    ('launch', false),
  ];

  static const _questions = ['storage', 'consent'];

  List<Map<String, Object?>> get _problems => [
    for (final (key, required) in _sections)
      if (required && statuses[key] != 'approved')
        {'section': key, 'code': 'section_not_approved', 'count': 1},
    if (!signed) {'section': 'launch', 'code': 'protocol_unsigned', 'count': 1},
  ];

  Map<String, Object?> _overview() => {
    'launched': launched,
    'launchedAt': null,
    'protocolVersion': signed ? 1 : 0,
    'signedBy': null,
    'signedAt': signed ? '2026-11-13T10:00:00Z' : null,
    'sections': [
      for (final (key, required) in _sections)
        {
          'key': key,
          'required': required,
          'status': statuses[key] ?? 'not_started',
          'assigneeId': assignees[key],
          'dueOn': null,
          'note': '',
          'settings': key == 'privacy' ? ['core.self_deletion'] : <String>[],
          'updatedAt': null,
        },
    ],
    'legal': [
      for (final q in _questions)
        {
          'question': q,
          'title': q,
          'answer': answers[q]?['answer'] ?? '',
          'basis': answers[q]?['basis'] ?? '',
          'answeredBy': null,
          'answeredAt': null,
          'approvedBy': null,
          'approvedAt': answers[q]?['approved'] == true
              ? '2026-11-13T10:00:00Z'
              : null,
        },
    ],
    'problems': launched ? <Object?>[] : _problems,
  };

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    calls.add('${request.method} ${segments.join('/')}');
    final body = request.body.isEmpty
        ? const <String, Object?>{}
        : jsonDecode(request.body) as Map<String, Object?>;
    switch ((request.method, segments)) {
      case ('GET', ['status']):
        return _json(200, {
          'launched': launched,
          'manages': !launched && manages,
          'assigned': <String>[],
        });
      case ('GET', []):
        return _json(200, _overview());
      case ('PUT', ['sections', final key, 'assignee']):
        assignees[key] = body['userId'] as String?;
      case ('PUT', ['sections', final key, 'settings', _]):
        statuses[key] = 'in_progress';
        return http.Response('', 204);
      case ('POST', ['sections', final key, 'submit']):
        statuses[key] = 'review';
      case ('POST', ['sections', final key, 'review']):
        statuses[key] = body['approve'] == true ? 'approved' : 'in_progress';
      case ('PUT', ['legal', final q]):
        answers[q] = {'answer': body['answer'], 'basis': body['basis']};
        signed = false;
      case ('POST', ['legal', final q, 'approve']):
        answers[q]?['approved'] = true;
      case ('GET', ['protocol']):
        return http.Response.bytes(utf8.encode('# Протокол'), 200);
      case ('POST', ['protocol', 'sign']):
        signed = true;
      case ('POST', ['launch']):
        if (_problems.isNotEmpty) {
          return _json(409, {'code': 'setup_not_ready'});
        }
        launched = true;
      default:
        return _json(404, {'code': 'not_found'});
    }
    return _json(200, _overview());
  }

  http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
