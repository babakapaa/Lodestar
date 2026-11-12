import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'admin_navigation.dart';
import 'app_harness.dart';

void main() {
  testWidgets('an administrator creates a role with permissions', (
    tester,
  ) async {
    final server = FakeServer()..permissions = ['core.administrator'];
    await pumpApp(tester, size: const Size(390, 1200), server: server);
    await openAdmin(tester, 'Роли');

    expect(find.text('Преподаватель'), findsOneWidget);
    expect(find.text('Есть у всех'), findsOneWidget);

    await tester.tap(find.text('Новая роль'));
    await tester.pumpAndSettle();
    await tester.enterText(find.byType(TextField), 'Куратор');
    await tester.tap(find.text('Состав групп'));
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();

    final created = server.directory.roles.last;
    expect(created['name'], 'Куратор');
    expect(created['permissions'], ['core.structure.members']);
    expect(find.text('Куратор'), findsOneWidget);
  });

  testWidgets('dragging a role sends its new position', (tester) async {
    final server = FakeServer()..permissions = ['core.administrator'];
    await pumpApp(tester, size: const Size(390, 1200), server: server);
    await openAdmin(tester, 'Роли');

    final handle = find.byIcon(Icons.drag_handle).last;
    final gesture = await tester.startGesture(tester.getCenter(handle));
    await gesture.moveBy(const Offset(0, -40));
    await tester.pump();
    await gesture.moveBy(const Offset(0, -80));
    await tester.pump();
    await gesture.up();
    await tester.pumpAndSettle();

    expect(server.directory.moves, [('r-teacher', 2)]);
  });

  testWidgets('a role is granted to a person and revoked', (tester) async {
    final server = FakeServer()..permissions = ['core.roles.manage'];
    await pumpApp(tester, size: const Size(390, 1200), server: server);
    await openAdmin(tester, 'Люди и роли');

    await tester.enterText(find.byType(TextField), 'ерлан');
    await tester.pump(const Duration(milliseconds: 300));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Ерлан Ким'));
    await tester.pumpAndSettle();

    expect(find.text('Ролей пока нет'), findsOneWidget);
    await tester.tap(find.text('Выдать роль'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Преподаватель'));
    await tester.pumpAndSettle();

    expect(server.directory.assignments.single['roleId'], 'r-teacher');
    expect(find.text('Преподаватель'), findsOneWidget);

    await tester.tap(find.byTooltip('Забрать роль'));
    await tester.pumpAndSettle();
    expect(server.directory.assignments, isEmpty);
    expect(find.text('Ролей пока нет'), findsOneWidget);
  });
}
