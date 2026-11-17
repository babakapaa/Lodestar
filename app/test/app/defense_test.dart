import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';
import 'search_map_test.dart' show scrollTo;
import 'study_test.dart' show openAssignment;

Map<String, Object?> slot(String id, Duration from) => {
  'id': id,
  'assignmentId': 'a1',
  'startsAt': DateTime.now().add(from).toUtc().toIso8601String(),
  'minutes': 15,
  'place': 'Ауд. 214',
  'capacity': 1,
};

void main() {
  testWidgets('a student books a defense time and sees question types', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.assignments.first['defense'] = 'full';
    server.study.defense.slots
      ..add(slot('slot1', const Duration(days: 1)))
      ..add(slot('slot2', const Duration(days: 2)));
    await pumpApp(tester, size: const Size(390, 2000), server: server);
    await openAssignment(tester);

    await scrollTo(tester, find.text('Записаться'));
    expect(find.text('Защита: полная защита'), findsOneWidget);
    expect(find.text('Вы ещё не записались.'), findsOneWidget);
    await tester.tap(find.text('Какие будут вопросы'));
    await tester.pumpAndSettle();
    expect(
      find.text('Как вы нашли главный источник и почему отбросили другие?'),
      findsOneWidget,
    );

    await tester.tap(find.byKey(const ValueKey('defense-book')));
    await tester.pumpAndSettle();
    expect(find.text('Запись на защиту'), findsOneWidget);
    await tester.tap(find.byKey(const ValueKey('slot-slot2')));
    await tester.pumpAndSettle();

    expect(server.study.defense.slotId, 'slot2');
    expect(find.textContaining('Вы записаны:'), findsOneWidget);
    expect(find.text('Выбрать другое время'), findsOneWidget);

    await tester.tap(find.text('Отменить запись'));
    await tester.pumpAndSettle();
    expect(server.study.defense.slotId, isNull);
    expect(find.text('Вы ещё не записались.'), findsOneWidget);
  });

  testWidgets('a teacher holds a defense with own questions and a grade', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    server.study.assignments.first['defense'] = 'full';
    await pumpApp(tester, size: const Size(390, 2000), server: server);
    await openAssignment(tester);

    await tester.tap(find.byKey(const ValueKey('defenses')));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('add-slots')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('slots-count')), '2');
    await tester.tap(find.text('Добавить'));
    await tester.pumpAndSettle();
    expect(server.study.defense.slots, hasLength(2));
    expect(find.text('Никто не записан'), findsNWidgets(2));

    await tester.tap(find.text('Ерлан Ким'));
    await tester.pumpAndSettle();
    expect(find.text('Вопросов пока нет.'), findsOneWidget);
    for (final text in [
      'Что такое мембрана клетки?',
      'Как вы нашли главный источник?',
    ]) {
      await tester.tap(find.byKey(const ValueKey('add-question')));
      await tester.pumpAndSettle();
      await tester.enterText(find.byKey(const ValueKey('question-text')), text);
      await tester.tap(find.text('Добавить'));
      await tester.pumpAndSettle();
    }
    expect(find.text('1. Что такое мембрана клетки?'), findsOneWidget);
    expect(find.text('2. Как вы нашли главный источник?'), findsOneWidget);

    await tester.tap(find.text('Понимает').first);
    await tester.enterText(
      find.byKey(const ValueKey('note-0')),
      'Ответил без подсказок',
    );
    await tester.tap(find.byTooltip('Убрать вопрос').last);
    await tester.pumpAndSettle();
    expect(find.text('2. Как вы нашли главный источник?'), findsNothing);

    await scrollTo(tester, find.byKey(const ValueKey('draft-protocol')));
    await tester.tap(find.byKey(const ValueKey('draft-protocol')));
    await tester.pumpAndSettle();
    expect(
      find.text('Студент уверенно объяснил строение клетки.'),
      findsOneWidget,
    );

    await scrollTo(tester, find.byKey(const ValueKey('finish-defense')));
    await tester.tap(find.byKey(const ValueKey('finish-defense')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('finish-points')), '9');
    await tester.tap(find.text('Завершить'));
    await tester.pumpAndSettle();

    final defense = server.study.defense;
    expect(defense.status, 'held');
    expect(defense.questions, [
      {
        'kind': 'understanding',
        'text': 'Что такое мембрана клетки?',
        'mark': 2,
        'note': 'Ответил без подсказок',
      },
    ]);
    expect(defense.summary, 'Студент уверенно объяснил строение клетки.');
    expect(server.study.grades['a1']?['points'], 9);
    expect(find.text('Защита завершена'), findsOneWidget);
    expect(find.text('Защита прошла'), findsOneWidget);
  });
}
