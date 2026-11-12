import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'admin_navigation.dart';
import 'app_harness.dart';

void main() {
  testWidgets('an administrator builds a group and adds a student', (
    tester,
  ) async {
    final server = FakeServer()..permissions = ['core.administrator'];
    await pumpApp(tester, size: const Size(390, 1200), server: server);
    await openAdmin(tester, 'Подразделения');

    await tester.tap(find.text('Факультет информатики'));
    await tester.pumpAndSettle();
    expect(find.text('Здесь пока пусто'), findsOneWidget);

    await tester.tap(find.text('Добавить'));
    await tester.pumpAndSettle();
    await tester.enterText(find.widgetWithText(TextField, 'Название'), 'ИС-21');
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();

    final group = server.directory.units.last;
    expect(group['parentId'], 'f1');
    expect(group['kind'], 'group');

    await tester.tap(find.text('ИС-21'));
    await tester.pumpAndSettle();
    expect(find.text('Факультет информатики / Группа'), findsOneWidget);

    await tester.tap(find.text('Добавить человека'));
    await tester.pumpAndSettle();
    await tester.enterText(find.byType(TextField), 'Ерлан');
    await tester.pump(const Duration(milliseconds: 300));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Ерлан Ким'));
    await tester.pumpAndSettle();

    expect(server.directory.members[group['id']], {'u2'});
    expect(find.text('erlan@kaznu.kz'), findsOneWidget);

    await tester.tap(find.byTooltip('Убрать'));
    await tester.pumpAndSettle();
    expect(server.directory.members[group['id']], isEmpty);
  });

  testWidgets('a teacher manages members but not the structure', (
    tester,
  ) async {
    final server = FakeServer()..permissions = ['core.structure.members'];
    await pumpApp(tester, size: const Size(390, 1200), server: server);
    await openAdmin(tester, 'Подразделения');

    expect(find.text('Добавить'), findsNothing);
    await tester.tap(find.text('Факультет информатики'));
    await tester.pumpAndSettle();
    expect(find.byIcon(Icons.more_vert), findsNothing);
    expect(find.text('Добавить человека'), findsOneWidget);
  });

  testWidgets('a unit goes to the archive and comes back', (tester) async {
    final server = FakeServer()..permissions = ['core.structure.manage'];
    await pumpApp(tester, size: const Size(390, 1200), server: server);
    await openAdmin(tester, 'Подразделения');

    await tester.tap(find.text('Факультет информатики'));
    await tester.pumpAndSettle();
    await tester.tap(find.byIcon(Icons.more_vert));
    await tester.pumpAndSettle();
    await tester.tap(find.text('В архив'));
    await tester.pumpAndSettle();
    expect(server.directory.units.first['archivedAt'], isNotNull);
    expect(find.text('В архиве'), findsOneWidget);

    await tester.tap(find.byIcon(Icons.more_vert));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Вернуть из архива'));
    await tester.pumpAndSettle();
    expect(server.directory.units.first['archivedAt'], isNull);
  });
}
