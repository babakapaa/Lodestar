import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';
import 'search_map_test.dart' show scrollTo;

Future<void> openInsights(WidgetTester tester) async {
  await tester.tap(find.text('Учёба').first);
  await tester.pumpAndSettle();
  await tester.tap(find.text('Биология'));
  await tester.pumpAndSettle();
  await scrollTo(tester, find.byKey(const ValueKey('course-insights')));
  await tester.tap(find.byKey(const ValueKey('course-insights')));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('topics stay hidden below five students', (tester) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    await pumpApp(tester, size: const Size(390, 2400), server: server);
    await openInsights(tester);

    expect(
      find.text(
        'Навигатор спрашивали меньше 5 студентов, поэтому темы пока скрыты.',
      ),
      findsOneWidget,
    );
  });

  testWidgets('a teacher sees topics and a suggestion for the next class', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    final topic = {
      'materialId': 'm1',
      'title': 'Строение клетки',
      'students': 7,
      'recent': 6,
    };
    server.navigator.topics = {
      'threshold': 5,
      'weeks': 8,
      'students': 9,
      'topics': [topic],
      'intents': [
        {'intent': 'understand', 'students': 6},
      ],
      'trend': [
        {'week': '2026-09-07', 'students': null},
        {'week': '2026-09-14', 'students': 8},
      ],
      'suggestion': topic,
    };
    await pumpApp(tester, size: const Size(390, 2400), server: server);
    await openInsights(tester);

    expect(find.text('Спрашивали 9 студентов'), findsOneWidget);
    expect(find.byKey(const ValueKey('insights-suggestion')), findsOneWidget);
    expect(find.textContaining('«Строение клетки»'), findsOneWidget);
    expect(find.byKey(const ValueKey('topic-m1')), findsOneWidget);
    expect(find.text('меньше 5'), findsOneWidget);
    expect(find.text('8 студ.'), findsOneWidget);
  });
}
