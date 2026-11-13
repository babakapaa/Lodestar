import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'admin_navigation.dart';
import 'app_harness.dart';
import 'search_map_test.dart' show scrollTo;

Future<void> openLeague(WidgetTester tester) async {
  await tester.tap(find.text('Учёба').first);
  await tester.pumpAndSettle();
  await tester.tap(find.text('Биология'));
  await tester.pumpAndSettle();
  await scrollTo(tester, find.byKey(const ValueKey('course-league')));
  await tester.tap(find.byKey(const ValueKey('course-league')));
  await tester.pumpAndSettle();
}

Map<String, Object?> task(String id) => {
  'id': id,
  'leagueId': 'l1',
  'title': 'Бот для библиотеки',
  'about': 'Найти статьи быстрее',
  'partner': 'Библиотека',
  'points': 60,
  'minSize': 2,
  'maxSize': 3,
  'selfFormed': true,
  'dueAt': null,
  'closedAt': null,
  'createdAt': '2026-11-13T10:00:00Z',
};

void main() {
  testWidgets(
    'a student sees the table, hides behind a nickname and disputes',
    (tester) async {
      final server = FakeServer();
      await pumpApp(tester, size: const Size(390, 2400), server: server);
      await openLeague(tester);

      expect(find.text('Сезон 2'), findsOneWidget);
      expect(find.text('Бронза'), findsOneWidget);
      expect(find.text('До дивизиона «Серебро»: 30'), findsOneWidget);
      expect(find.text('Айгерим Сапарова (вы)'), findsOneWidget);
      expect(
        find.text('рост +50 к прошлому сезону, лучший рост'),
        findsOneWidget,
      );

      await tester.tap(find.byKey(const ValueKey('league-display')));
      await tester.pumpAndSettle();
      await tester.tap(find.text('под ником'));
      await tester.pumpAndSettle();
      await tester.enterText(
        find.byKey(const ValueKey('league-nickname')),
        'Клетка',
      );
      await tester.tap(find.text('Сохранить'));
      await tester.pumpAndSettle();
      expect(server.league.standing['display'], 'nickname');
      expect(server.league.standing['nickname'], 'Клетка');
      expect(find.text('В таблице: под ником'), findsOneWidget);

      await scrollTo(tester, find.byKey(const ValueKey('dispute-e1')));
      expect(find.text('8 из 10, вес 100'), findsOneWidget);
      await tester.tap(find.byKey(const ValueKey('dispute-e1')));
      await tester.pumpAndSettle();
      await tester.enterText(
        find.byKey(const ValueKey('team-text')),
        'Оценку исправили',
      );
      await tester.tap(find.text('Отправить'));
      await tester.pumpAndSettle();
      expect(server.league.disputes.single['entryId'], 'e1');
      expect(find.text('Оценку исправили'), findsNothing);
      expect(find.textContaining('Оценку исправили'), findsOneWidget);
    },
  );

  testWidgets('a teacher confirms a move, adjusts points and sets rules', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    await pumpApp(tester, size: const Size(390, 2400), server: server);
    await openLeague(tester);

    expect(find.text('Переходы на подтверждение'), findsOneWidget);
    await tester.tap(find.byKey(const ValueKey('move-approve-m1')));
    await tester.pumpAndSettle();
    expect(server.league.pending, isEmpty);
    expect(find.text('Переходы на подтверждение'), findsNothing);

    await tester.tap(find.byKey(const ValueKey('player-u2')));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('league-adjust')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('adjust-points')), '15');
    await tester.enterText(
      find.byKey(const ValueKey('adjust-reason')),
      'Олимпиада',
    );
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();
    expect(server.league.adjustments.single['points'], 15);
    expect(find.text('поправка: Олимпиада'), findsOneWidget);
    await tester.tap(find.byTooltip('Назад'));
    await tester.pumpAndSettle();

    await tester.tap(find.byKey(const ValueKey('league-rules')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('weight-a1')), '50');
    await tester.tap(find.byKey(const ValueKey('division-add')));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('division-name-2')),
      'Золото',
    );
    await tester.tap(find.byKey(const ValueKey('league-confirm')));
    await scrollTo(tester, find.byKey(const ValueKey('rules-save')));
    await tester.tap(find.byKey(const ValueKey('rules-save')));
    await tester.pumpAndSettle();

    expect(server.league.weight, 50);
    expect(server.league.confirm, isFalse);
    expect(server.league.divisions.last, {
      'name': 'Золото',
      'from': 300,
      'gpuHours': 0,
    });
  });

  testWidgets('a teacher publishes a task and reviews a similar solution', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    await pumpApp(tester, size: const Size(390, 2400), server: server);
    await openLeague(tester);

    await tester.tap(find.byKey(const ValueKey('task-new')));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('task-title')),
      'Бот для библиотеки',
    );
    await tester.enterText(find.byKey(const ValueKey('task-points')), '60');
    await tester.tap(find.byKey(const ValueKey('task-save')));
    await tester.pumpAndSettle();
    expect(server.league.tasks.single['points'], 60);

    await tester.tap(find.text('Бот для библиотеки'));
    await tester.pumpAndSettle();
    expect(
      find.text(
        'Похоже на решение команды «Старички» по задаче «Бот прошлого года» (93%)',
      ),
      findsOneWidget,
    );
    expect(find.text('Ерлан Ким, Мадина Омарова'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('review-sol1')));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('passed-u4')));
    await tester.tap(find.byKey(const ValueKey('review-save')));
    await tester.pumpAndSettle();
    expect(server.league.review, {
      'points': 60,
      'passed': ['u2'],
      'comment': '',
    });
    expect(find.text('Оценено: 60, защитили: Ерлан Ким'), findsOneWidget);
  });

  testWidgets('a student opens a league task and looks for a team', (
    tester,
  ) async {
    final server = FakeServer();
    server.league.tasks.add(task('task1'));
    await pumpApp(tester, size: const Size(390, 2400), server: server);
    await openLeague(tester);

    await tester.tap(find.byKey(const ValueKey('task-task1')));
    await tester.pumpAndSettle();
    expect(find.text('Найти статьи быстрее'), findsOneWidget);
    expect(find.textContaining('от Библиотека'), findsOneWidget);
  });

  testWidgets('an administrator adds an award', (tester) async {
    final server = FakeServer()..permissions = ['league.awards.manage'];
    await pumpApp(tester, size: const Size(390, 1400), server: server);
    await openAdmin(tester, 'Награды лиги');

    await tester.tap(find.byKey(const ValueKey('award-new')));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('award-title')),
      'Лучший рост',
    );
    await tester.enterText(find.byKey(const ValueKey('award-gpu')), '20');
    await tester.enterText(find.byKey(const ValueKey('award-quantity')), '3');
    await tester.tap(find.byKey(const ValueKey('award-save')));
    await tester.pumpAndSettle();

    expect(server.league.awards.single['quantity'], 3);
    expect(find.text('20 ч GPU, вручено 0 из 3'), findsOneWidget);
  });
}
