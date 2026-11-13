import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';
import 'study_test.dart' show goBack, openAssignment;

Map<String, Object?> source(String id, String kind, String title) => {
  'id': id,
  'kind': kind,
  'title': title,
  'authors': ['Иванов, А. Б.'],
  'year': '2021',
  'container': kind == 'article' ? 'Вестник биологии' : '',
  'publisher': '',
  'city': '',
  'volume': '',
  'issue': '',
  'pages': '',
  'url': kind == 'web' ? 'https://cells.example.org/a' : '',
  'doi': '',
  'isbn': '',
  'accessedAt': '',
  'note': '',
};

Future<void> scrollTo(WidgetTester tester, Finder finder) =>
    tester.dragUntilVisible(
      finder,
      find.byType(ListView).last,
      const Offset(0, -300),
    );

void main() {
  testWidgets('a student builds a search map and hands it in', (tester) async {
    final server = FakeServer();
    server.study.assignments.first['searchMap'] = {
      'minSources': 1,
      'minScholarly': 50,
    };
    server.sources.sources['s1'] = source('s1', 'article', 'Строение клетки');
    server.sources.sources['s2'] = source('s2', 'web', 'Клетки для детей');
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: server,
    );
    await openAssignment(tester);

    expect(
      find.text('Карта поиска ещё не готова, сдать работу нельзя'),
      findsOneWidget,
    );
    await tester.tap(find.text('Собрать карту поиска'));
    await tester.pumpAndSettle();
    expect(find.text('Опишите путь поиска'), findsOneWidget);

    await tester.tap(find.text('Добавить из моих источников'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('В работу').first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Добавить из моих источников'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Отброшен'));
    await tester.pumpAndSettle();

    await tester.enterText(
      find.widgetWithText(TextField, 'Почему выбран'),
      'Свежие данные',
    );
    await tester.pumpAndSettle();
    expect(
      find.text('Напишите конкретнее: что именно вы проверили или взяли'),
      findsOneWidget,
    );
    await tester.enterText(
      find.widgetWithText(TextField, 'Почему выбран'),
      'Свежие данные о мембране с электронной микроскопии',
    );
    await tester.enterText(
      find.widgetWithText(TextField, 'Почему отброшен'),
      'Упрощено для школьников, нет ссылок на исследования',
    );
    await scrollTo(tester, find.text('Сохранить карту'));
    await tester.enterText(
      find.widgetWithText(
        TextField,
        'Главные повороты: сначала искал про X, понял, что важнее Y',
      ),
      'Сначала искал про ядро, потом понял, что важнее мембрана.',
    );
    await tester.pumpAndSettle();
    await scrollTo(tester, find.text('Сохранить карту'));
    await tester.tap(find.text('Сохранить карту'));
    await tester.pumpAndSettle();

    final saved = server.study.maps['a1']!;
    expect(saved['used'], hasLength(1));
    expect(saved['rejected'], hasLength(1));
    expect(
      (saved['used']! as List<Object?>).single,
      containsPair('reference', contains('Строение клетки')),
    );

    await tester.tap(find.byTooltip('Выгрузить в .md'));
    await tester.pumpAndSettle();
    final markdown = utf8.decode(app.saver.saved['search-map.md']!);
    expect(markdown, contains('## Отброшенные источники'));

    await goBack(tester);
    expect(find.text('Карта поиска приложится к работе'), findsOneWidget);
    await tester.enterText(find.byType(TextField), 'Клетка и мембрана');
    await tester.tap(find.text('Сдать работу'));
    await tester.pumpAndSettle();
    final submission = server.study.submissions.single;
    expect(submission['searchMap'], isNotNull);
    expect(find.text('Путь поиска'), findsOneWidget);
  });

  testWidgets('a teacher asks for a search map and reads it', (tester) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    server.study.submissions.add({
      'id': 's1',
      'assignmentId': 'a1',
      'clientId': 'x',
      'body': 'Ответ',
      'submittedAt': '2026-11-11T10:00:00Z',
      'files': <Object?>[],
      'searchMap': {
        'used': [
          {
            'sourceId': 's1',
            'reference': 'Иванов, А. Б. Строение клетки',
            'kind': 'article',
            'doi': '',
            'url': '',
            'found': 'eLibrary',
            'why': 'Свежие данные',
            'taken': 'Схема мембраны',
          },
        ],
        'rejected': <Object?>[],
        'path': 'Сначала ядро, потом мембрана',
        'queries': ['строение мембраны'],
      },
    });
    await pumpApp(tester, size: const Size(390, 1400), server: server);
    await openAssignment(tester);
    await tester.tap(find.byTooltip('Изменить задание'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Нужна карта поиска'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.widgetWithText(TextField, 'Минимум источников'),
      '3',
    );
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();
    expect(server.study.assignments.first['searchMap'], {
      'minSources': 3,
      'minScholarly': 50,
    });

    await tester.tap(find.text('Ерлан Ким'));
    await tester.pumpAndSettle();
    expect(find.text('1. Иванов, А. Б. Строение клетки'), findsOneWidget);
    expect(find.text('Где найден: eLibrary'), findsOneWidget);
    expect(find.text('Сначала ядро, потом мембрана'), findsOneWidget);
  });
}
