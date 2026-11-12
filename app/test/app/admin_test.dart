import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';

Future<void> openProfile(WidgetTester tester) async {
  await tester.tap(find.byTooltip('Профиль и настройки'));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('a student does not see administration', (tester) async {
    await pumpApp(tester, size: const Size(390, 1000));
    await openProfile(tester);

    expect(find.text('Администрирование'), findsNothing);
  });

  testWidgets('an administrator changes university settings and modules', (
    tester,
  ) async {
    final server = FakeServer()..permissions = ['core.administrator'];
    await pumpApp(tester, size: const Size(390, 1200), server: server);
    await openProfile(tester);

    await tester.tap(find.text('Администрирование'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Настройки вуза'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Люди сами удаляют свой аккаунт'));
    await tester.pumpAndSettle();
    expect(server.settings['core.self_deletion'], isFalse);

    await tester.tap(find.text('Дней до удаления аккаунта'));
    await tester.pumpAndSettle();
    expect(find.text('От 0 до 90'), findsOneWidget);
    await tester.enterText(find.byType(TextField), '30');
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();
    expect(server.settings['core.deletion_grace_days'], 30);

    await tester.tap(find.text('study'));
    await tester.pumpAndSettle();
    expect(server.modules.first['enabled'], isFalse);
  });

  testWidgets('the audit log names who changed what', (tester) async {
    final server = FakeServer()..permissions = ['core.audit.view'];
    await pumpApp(tester, size: const Size(390, 1000), server: server);
    await openProfile(tester);

    await tester.tap(find.text('Администрирование'));
    await tester.pumpAndSettle();
    expect(find.text('Настройки вуза'), findsNothing);
    await tester.tap(find.text('Журнал действий'));
    await tester.pumpAndSettle();

    expect(find.text('Изменена настройка: core.self_deletion'), findsOneWidget);
    expect(find.textContaining('Дана Ахметова'), findsOneWidget);
  });
}
