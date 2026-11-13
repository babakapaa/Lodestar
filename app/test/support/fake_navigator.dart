import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeNavigator {
  bool unavailable = false;
  String support = '';
  final List<Map<String, Object?>> journal = [];
  final Map<String, Map<String, Object?>> dialogs = {};
  final Map<String, List<Map<String, Object?>>> messages = {};
  final List<Map<String, Object?>> questions = [];
  final List<Map<String, Object?>> materials = [];
  int _next = 1;

  static const _steps = [
    'question',
    'direction',
    'source',
    'fragment',
    'example',
  ];

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    return switch ((request.method, segments)) {
      ('GET', ['dialogs']) => _json(200, dialogs.values.toList()),
      ('POST', ['dialogs']) => _create(request),
      ('GET', ['dialogs', final id]) => _details(id),
      ('DELETE', ['dialogs', final id]) => _delete(id),
      ('POST', ['dialogs', final id, 'messages']) => _ask(request, id),
      ('GET', ['journal']) => _json(200, [
        for (final e in journal)
          if (request.url.queryParameters['scopeId'] == null ||
              e['scopeId'] == request.url.queryParameters['scopeId'])
            e,
      ]),
      ('DELETE', ['journal', final id]) => _forget(id),
      ('GET', ['support']) => _json(200, {'contacts': support}),
      ('GET', ['courses', _, 'materials']) => _json(200, materials),
      ('POST', ['courses', final course, 'materials']) => _addMaterial(
        request,
        course,
      ),
      ('DELETE', ['materials', final id]) => _deleteMaterial(id),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  http.Response _create(http.Request request) {
    final body = _body(request);
    final id = body['id']! as String;
    final created = !dialogs.containsKey(id);
    dialogs.putIfAbsent(
      id,
      () => {
        'id': id,
        'scopeKind': body['scopeKind'] ?? '',
        'scopeId': body['scopeId'] ?? '',
        'title': '',
        'createdAt': '2026-11-13T10:00:00Z',
        'changedAt': '2026-11-13T10:00:00Z',
      },
    );
    messages.putIfAbsent(id, () => []);
    return _json(created ? 201 : 200, _full(id));
  }

  http.Response _details(String id) {
    if (!dialogs.containsKey(id)) {
      return _json(404, {'code': 'dialog_not_found'});
    }
    return _json(200, _full(id));
  }

  http.Response _delete(String id) {
    dialogs.remove(id);
    messages.remove(id);
    return http.Response('', 204);
  }

  http.Response _ask(http.Request request, String id) {
    if (unavailable) {
      return _json(503, {'code': 'navigator_unavailable'});
    }
    final body = _body(request);
    questions.add(body);
    final list = messages[id]!;
    final more = body['more'] == true;
    final previous = list.isEmpty ? null : list.last;
    final step = more && previous != null
        ? _steps[_steps.indexOf(previous['step']! as String) + 1]
        : 'question';
    final intent = body['intent'] ?? 'understand';
    final text = body['text'] as String? ?? '';
    final question = _message(list.length + 1, 'student', text, intent, step);
    final reply =
        _message(
          list.length + 2,
          'navigator',
          'Мембрана «отделяет содержимое клетки» [1]. Что вы уже знаете?',
          intent,
          step,
        )..addAll({
          'support': text.contains('тяжело')
              ? 'Психологическая служба: 112'
              : '',
          'citations': [
            {
              'n': 1,
              'materialId': 'm1',
              'title': 'Строение клетки',
              'reference': 'Иванов И. И. Биология. 2020.',
              'url': '',
              'text': 'Мембрана отделяет содержимое клетки от внешней среды.',
            },
          ],
        });
    list
      ..add(question)
      ..add(reply);
    final dialog = dialogs[id]!;
    if (text.isNotEmpty) {
      journal.add({
        'id': 'j${_next++}',
        'scopeKind': dialog['scopeKind'],
        'scopeId': dialog['scopeId'],
        'query': text,
        'sources': [
          {'title': 'Строение клетки', 'reference': '', 'url': ''},
        ],
        'createdAt': '2026-11-13T10:00:00Z',
      });
    }
    if ((dialog['title']! as String).isEmpty) {
      dialog['title'] = text;
    }
    return _json(201, {'question': question, 'reply': reply});
  }

  http.Response _forget(String id) {
    journal.removeWhere((e) => e['id'] == id);
    return http.Response('', 204);
  }

  http.Response _addMaterial(http.Request request, String course) {
    final body = _body(request);
    final material = {
      'id': 'm${_next++}',
      'courseId': course,
      'title': body['title'],
      'reference': body['reference'] ?? '',
      'url': body['url'] ?? '',
      'passages': 2,
      'indexed': 0,
      'createdAt': '2026-11-13T10:00:00Z',
    };
    materials.add(material);
    return _json(201, material);
  }

  http.Response _deleteMaterial(String id) {
    materials.removeWhere((m) => m['id'] == id);
    return http.Response('', 204);
  }

  Map<String, Object?> _full(String id) => {
    ...dialogs[id]!,
    'brief': dialogs[id]!['scopeKind'] == 'assignment'
        ? {
            'kind': 'assignment',
            'id': dialogs[id]!['scopeId'],
            'title': 'Эссе о клетке',
            'description': '',
            'course': 'Биология',
            'courseId': 'c1',
            'teaches': false,
            'dueAt': null,
            'mode': 'study',
          }
        : null,
    'messages': messages[id],
  };

  static Map<String, Object?> _message(
    int seq,
    String role,
    String text,
    Object? intent,
    String step,
  ) => {
    'seq': seq,
    'role': role,
    'text': text,
    'intent': intent,
    'step': role == 'navigator' ? step : '',
    'lastStep': role == 'navigator' && step == 'example',
    'support': '',
    'citations': <Object?>[],
    'createdAt': '2026-11-13T10:00:00Z',
  };

  static Map<String, Object?> _body(http.Request request) =>
      jsonDecode(request.body) as Map<String, Object?>;

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
