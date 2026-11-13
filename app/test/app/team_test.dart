import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';
import 'search_map_test.dart' show scrollTo;
import 'study_test.dart' show goBack, openAssignment;

const _rules = {'minSize': 2, 'maxSize': 3, 'selfFormed': true};

void main() {
  testWidgets('a student builds a team and proposes the work', (tester) async {
    final server = FakeServer();
    server.study.assignments.first['team'] = _rules;
    server.teams.seeking['u2'] = 'Люблю микроскоп';
    await pumpApp(tester, size: const Size(390, 2000), server: server);
    await openAssignment(tester);

    expect(find.text('Найти команду'), findsOneWidget);
    await tester.tap(find.byKey(const ValueKey('team-find')));
    await tester.pumpAndSettle();
    expect(find.text('Люблю микроскоп'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('team-create')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('team-text')), 'Клетки');
    await tester.tap(find.text('Создать команду').last);
    await tester.pumpAndSettle();

    expect(server.teams.teams.single['name'], 'Клетки');
    expect(find.text('Состав: 1'), findsOneWidget);
    expect(find.text('JOIN-t1'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('team-invite')));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('pick-u2')));
    await tester.pumpAndSettle();
    expect(server.teams.invited['t1'], ['u2']);
    expect(find.byKey(const ValueKey('invited-u2')), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('team-label')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('team-text')), 'Автор');
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();

    await goBack(tester);
    await goBack(tester);
    await scrollTo(tester, find.byKey(const ValueKey('team-propose')));
    await tester.enterText(
      find.byKey(const ValueKey('team-work')),
      'Клетка состоит из ядра',
    );
    await tester.tap(find.byKey(const ValueKey('team-propose')));
    await tester.pumpAndSettle();

    final proposal = server.teams.proposal!;
    expect((proposal['work']! as Map)['body'], 'Клетка состоит из ядра');
    expect(find.text('Вариант на согласовании'), findsOneWidget);
    expect(find.text('Предложить новый вариант'), findsOneWidget);
  });

  testWidgets('a teacher places students and removes one', (tester) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    server.study.assignments.first['team'] = _rules;
    await pumpApp(tester, size: const Size(390, 2000), server: server);
    await openAssignment(tester);

    await tester.tap(find.byKey(const ValueKey('teams')));
    await tester.pumpAndSettle();
    expect(find.text('Без команды: 3'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('team-distribute')));
    await tester.pumpAndSettle();
    expect(server.teams.teams, hasLength(1));
    expect(find.text('Команды: 1'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('team-t1')));
    await tester.pumpAndSettle();
    expect(find.text('Состав: 3'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('member-menu-u4')));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Убрать из команды'));
    await tester.pumpAndSettle();
    expect(server.teams.teams.single['members'], ['u1', 'u2']);
    expect(find.text('Состав: 2'), findsOneWidget);
  });

  testWidgets('a student fills the matching card', (tester) async {
    final server = FakeServer();
    await pumpApp(tester, size: const Size(390, 2000), server: server);
    await tester.tap(find.text('Учёба').first);
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('my-teams')));
    await tester.pumpAndSettle();
    expect(find.text('Команд пока нет'), findsOneWidget);

    await tester.tap(find.text('Карточка для подбора'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('profile-interests')),
      'Микробиология',
    );
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();
    expect(server.teams.card['interests'], 'Микробиология');
  });
}
