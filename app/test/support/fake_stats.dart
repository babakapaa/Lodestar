import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeStats {
  final List<Map<String, String>> queries = [];
  final List<Map<String, String>> entries = [];

  http.Response entriesResponse() => http.Response(
    jsonEncode(entries),
    200,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );

  http.Response handle(http.Request request) {
    final query = request.url.queryParameters;
    queries.add(query);
    final parent = query['parent'];
    final units = parent == null
        ? [
            {
              'unitId': 'bio',
              'kind': 'faculty',
              'name': 'Биофак',
              'people': 120,
              'hidden': false,
              'metrics': [
                {'key': 'study.submitted', 'value': 0.82},
                {'key': 'attendance.rate', 'value': 0.9},
                {'key': 'sandbox.hours', 'value': 12.5},
              ],
            },
            {
              'unitId': 'arts',
              'kind': 'faculty',
              'name': 'Худграф',
              'people': null,
              'hidden': true,
              'metrics': <Object?>[],
            },
          ]
        : [
            {
              'unitId': 'gen',
              'kind': 'department',
              'name': 'Кафедра генетики',
              'people': 14,
              'hidden': false,
              'metrics': <Object?>[],
            },
          ];
    return http.Response(
      jsonEncode({
        'parentId': parent,
        'days': int.parse(query['days'] ?? '30'),
        'threshold': 5,
        'full': false,
        'units': units,
      }),
      200,
      headers: {'content-type': 'application/json; charset=utf-8'},
    );
  }
}
