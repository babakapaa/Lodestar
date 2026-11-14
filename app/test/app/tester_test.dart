import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';

const _people = [
  {
    'key': 'demo-student',
    'name': 'Айгерим Сапарова',
    'title': 'Студентка',
    'email': 'aigerim@kaznu.kz',
  },
  {
    'key': 'demo-teacher',
    'name': 'Ерлан Жумабаев',
    'title': 'Преподаватель',
    'email': 'teacher@kaznu.kz',
  },
];

FakeServer _testerServer() =>
    FakeServer()..tester = {'tenant': 'kaznu', 'people': _people};

void main() {
  testWidgets('tester mode signs in without the browser', (tester) async {
    final app = await pumpApp(
      tester,
      size: const Size(390, 1000),
      signedIn: false,
      server: _testerServer(),
    );

    expect(find.text('Режим тестировщика'), findsOneWidget);
    await tester.tap(find.text('Ерлан Жумабаев, Преподаватель'));
    await tester.pumpAndSettle();

    final signIn = app.server.requests.lastWhere(
      (request) => request.url.path == '/v1/tester/sign-in',
    );
    expect(jsonDecode(signIn.body), containsPair('person', 'demo-teacher'));
    expect(find.byType(NavigationBar), findsOneWidget);
  });

  testWidgets('tester mode switches the user from the profile', (tester) async {
    final app = await pumpApp(
      tester,
      size: const Size(390, 1600),
      server: _testerServer(),
    );

    await tester.tap(find.byTooltip('Профиль и настройки'));
    await tester.pumpAndSettle();
    await tester.scrollUntilVisible(
      find.text('Ерлан Жумабаев, Преподаватель'),
      200,
    );
    final current = tester.widget<TextButton>(
      find.ancestor(
        of: find.text('Айгерим Сапарова, Студентка'),
        matching: find.byWidgetPredicate((widget) => widget is TextButton),
      ),
    );
    expect(current.onPressed, isNull);

    await tester.tap(find.text('Ерлан Жумабаев, Преподаватель'));
    await tester.pumpAndSettle();

    final paths = app.server.requests.map((request) => request.url.path);
    expect(
      paths,
      containsAllInOrder(['/v1/auth/logout', '/v1/tester/sign-in']),
    );
    expect(app.server.validAccess, 'access-2');
    expect(find.byType(NavigationBar), findsOneWidget);
  });

  testWidgets('a usual server shows no tester mode', (tester) async {
    await pumpApp(tester, size: const Size(390, 844), signedIn: false);

    expect(find.text('Войти через Microsoft'), findsOneWidget);
    expect(find.text('Режим тестировщика'), findsNothing);
  });
}
