import 'dart:convert';
import 'dart:typed_data';

import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:lodestar/core/files/opened_file.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'admin_navigation.dart';
import 'app_harness.dart';
import 'search_map_test.dart' show scrollTo;

Future<void> go(WidgetTester tester, String location) async {
  GoRouter.of(tester.element(find.byType(Scaffold).first)).go(location);
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('an administrator connects and disconnects an LMS', (
    tester,
  ) async {
    final server = FakeServer()..permissions = ['core.integrations.manage'];
    await pumpApp(tester, size: const Size(390, 1600), server: server);
    await openAdmin(tester, 'LMS');

    expect(
      find.text('https://campus.example.edu/v1/lti/tenants/kaznu/login'),
      findsOneWidget,
    );
    expect(find.text('LMS пока не подключены.'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('lms-new')));
    await tester.pumpAndSettle();
    final fields = {
      'name': 'Moodle',
      'issuer': 'https://moodle.kaznu.kz',
      'client': 'abc',
      'auth': 'https://moodle.kaznu.kz/auth',
      'token': 'https://moodle.kaznu.kz/token',
      'jwks': 'https://moodle.kaznu.kz/certs',
      'deployments': '1, 2',
    };
    for (final MapEntry(:key, :value) in fields.entries) {
      await tester.enterText(find.byKey(ValueKey('lms-$key')), value);
    }
    await tester.tap(find.byKey(const ValueKey('lms-save')));
    await tester.pumpAndSettle();

    expect(server.lti.platforms.single['deploymentIds'], ['1', '2']);
    expect(find.byKey(const ValueKey('lms-p1')), findsOneWidget);

    await tester.tap(find.byTooltip('Отключить LMS?'));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('lms-delete-confirm')));
    await tester.pumpAndSettle();
    expect(server.lti.platforms, isEmpty);
    expect(find.text('LMS пока не подключены.'), findsOneWidget);
  });

  testWidgets('a launch from the LMS signs in and opens the course', (
    tester,
  ) async {
    final server = FakeServer();
    await pumpApp(tester, size: const Size(390, 1600), server: server);

    await go(
      tester,
      '/lti?server=https://campus.example.edu&tenant=kaznu'
      '&ticket=lst_ticket&next=/study/courses/c1',
    );

    expect(server.validAccess, 'access-2');
    expect(find.text('Основы клеточной биологии'), findsOneWidget);
  });

  testWidgets('a launch for another university asks first', (tester) async {
    final server = FakeServer();
    await pumpApp(tester, size: const Size(390, 1600), server: server);

    await go(
      tester,
      '/lti?server=https://other.example.edu&tenant=kaznu&ticket=lst_ticket',
    );

    expect(server.validAccess, 'access-1');
    expect(find.byKey(const ValueKey('lti-enter')), findsOneWidget);
  });

  testWidgets('a broken launch offers the way home', (tester) async {
    await pumpApp(tester, size: const Size(390, 1600));

    await go(tester, '/lti?error=launch_rejected');

    expect(
      find.text('Не удалось войти из LMS. Откройте курс в LMS ещё раз.'),
      findsOneWidget,
    );
    await tester.tap(find.text('На главную'));
    await tester.pumpAndSettle();
    expect(find.byKey(const ValueKey('lti-enter')), findsNothing);
  });

  testWidgets('a teacher links an LMS course, picks the source and syncs', (
    tester,
  ) async {
    final server = FakeServer();
    await pumpApp(tester, size: const Size(390, 1600), server: server);

    await go(tester, '/lti/contexts/x1');
    expect(find.text('Biology 101 (BIO)'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('lti-link-c1')));
    await tester.pumpAndSettle();
    expect(find.text('Связан с «Биология»'), findsOneWidget);

    await tester.tap(find.text('LMS').last);
    await tester.pumpAndSettle();
    expect(server.lti.contexts['x1']!['assignmentSource'], 'lms');
    expect(
      find.text('Задания приходят из LMS, оценки возвращаются туда.'),
      findsOneWidget,
    );

    await tester.tap(find.byKey(const ValueKey('lti-sync')));
    await tester.pumpAndSettle();
    expect(server.lti.syncs, 1);
    expect(
      find.text(
        'Людей в LMS: 3, добавлено: 2, заданий: 1, оценок отправлено: 1',
      ),
      findsOneWidget,
    );

    await tester.tap(find.byKey(const ValueKey('lti-unlink')));
    await tester.pumpAndSettle();
    expect(server.lti.contexts['x1']!['link'], isNull);
  });

  testWidgets('a teacher exports the gradebook and imports students', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    server.lti.contexts['x1']!['link'] = {
      'kind': 'study.course',
      'id': 'c1',
      'title': 'Биология',
      'path': '/study/courses/c1',
    };
    final app = await pumpApp(
      tester,
      size: const Size(390, 2400),
      server: server,
    );
    await tester.tap(find.text('Учёба').first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Биология'));
    await tester.pumpAndSettle();
    await scrollTo(tester, find.byKey(const ValueKey('course-lms')));
    await tester.tap(find.byKey(const ValueKey('course-lms')));
    await tester.pumpAndSettle();

    expect(find.text('Biology 101 (BIO)'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('gradebook-xlsx')));
    await tester.pumpAndSettle();
    final saved = app.saver.saved.entries.single;
    expect(saved.key, endsWith('.xlsx'));
    expect(utf8.decode(saved.value), contains('xlsx'));

    app.opener.next = [
      OpenedFile(
        'group.csv',
        Uint8List.fromList(utf8.encode('email\nerlan@kaznu.kz\nnobody@x.kz\n')),
      ),
    ];
    await tester.tap(find.byKey(const ValueKey('import-students')));
    await tester.pumpAndSettle();
    expect(find.text('Ничего не добавлено. Ошибок: 1'), findsOneWidget);
    expect(
      find.text('Строка 3: такого человека нет в Lodestar'),
      findsOneWidget,
    );
    expect(server.study.imported, isEmpty);

    app.opener.next = [
      OpenedFile(
        'group.csv',
        Uint8List.fromList(utf8.encode('email\nerlan@kaznu.kz\n')),
      ),
    ];
    await tester.tap(find.byKey(const ValueKey('import-students')));
    await tester.pumpAndSettle();
    expect(find.text('Добавлено: 1, уже были в курсе: 0'), findsOneWidget);
    expect(server.study.imported, ['erlan@kaznu.kz']);
  });
}
