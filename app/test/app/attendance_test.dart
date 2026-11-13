import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';
import 'sandbox_test.dart' show go;
import 'search_map_test.dart' show scrollTo;

FakeServer attendanceServer() {
  final server = FakeServer();
  server.modules.add({
    'name': 'attendance',
    'version': '0.1.0',
    'optional': true,
    'enabled': true,
  });
  return server;
}

void main() {
  testWidgets('a teacher opens check-in, marks and closes it', (tester) async {
    final server = attendanceServer();
    server.study.role = 'teacher';
    final app = await pumpApp(
      tester,
      size: const Size(390, 2400),
      server: server,
    );
    await tester.tap(find.text('Учёба').first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Биология'));
    await tester.pumpAndSettle();
    await scrollTo(tester, find.byKey(const ValueKey('course-attendance')));
    await tester.tap(find.byKey(const ValueKey('course-attendance')));
    await tester.pumpAndSettle();

    expect(find.text('Давно не были'), findsOneWidget);
    expect(find.byKey(const ValueKey('away-u2')), findsOneWidget);
    expect(find.textContaining('3 пропусков подряд'), findsWidgets);

    await tester.tap(find.byKey(const ValueKey('attendance-export')));
    await tester.pumpAndSettle();
    await tester.tap(find.text('CSV'));
    await tester.pumpAndSettle();
    expect(app.saver.saved.keys.single, endsWith('.csv'));

    await tester.tap(find.byKey(const ValueKey('lesson-ser1-2026-09-08')));
    await tester.pumpAndSettle();
    expect(find.text('Код: 4821'), findsOneWidget);
    expect(find.byKey(const ValueKey('attendance-qr')), findsOneWidget);
    expect(find.text('Отметились 0 из 2'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('mark-u1')));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Есть').last);
    await tester.pumpAndSettle();
    expect(server.attendance.marked.single, {
      'session': 'a1',
      'user': 'u1',
      'status': 'present',
    });
    expect(find.text('Отметились 1 из 2'), findsOneWidget);

    await scrollTo(tester, find.byKey(const ValueKey('attendance-close')));
    await tester.tap(find.byKey(const ValueKey('attendance-close')));
    await tester.pumpAndSettle();
    expect(find.text('Отметка закрыта'), findsOneWidget);
    expect(find.byKey(const ValueKey('attendance-qr')), findsNothing);
  });

  testWidgets('a student checks in with a code from Today', (tester) async {
    final server = attendanceServer();
    server.attendance.open.add({
      'sessionId': 'a1',
      'courseTitle': 'Биология',
      'room': '301',
      'startsAt': '2026-09-08T09:00:00Z',
      'endsAt': '2026-09-08T10:20:00Z',
      'status': null,
      'methods': {'qr': true, 'code': true},
    });
    await pumpApp(tester, size: const Size(390, 2400), server: server);

    expect(find.text('Отметка на занятии'), findsOneWidget);
    await tester.tap(find.text('Отметиться').first);
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('attendance-code-field')),
      '1111',
    );
    await tester.tap(find.byKey(const ValueKey('attendance-check')));
    await tester.pumpAndSettle();
    expect(find.text('Код неверный или уже устарел'), findsOneWidget);

    await tester.enterText(
      find.byKey(const ValueKey('attendance-code-field')),
      '4821',
    );
    await tester.tap(find.byKey(const ValueKey('attendance-check')));
    await tester.pumpAndSettle();
    expect(find.text('«Биология»: Есть'), findsOneWidget);
    expect(server.attendance.checks.last, {'code': '4821'});
  });

  testWidgets('a scanned QR link checks the student in at once', (
    tester,
  ) async {
    final server = attendanceServer();
    await pumpApp(tester, size: const Size(390, 2400), server: server);
    await go(tester, '/attend?s=a1&t=tok');

    expect(find.byKey(const ValueKey('attendance-done')), findsOneWidget);
    expect(server.attendance.checks.single, {
      'sessionId': 'a1',
      'token': 'tok',
    });
  });
}
