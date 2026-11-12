import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';

Future<void> openPrivacy(WidgetTester tester) async {
  await tester.tap(find.byTooltip('Профиль и настройки'));
  await tester.pumpAndSettle();
  await tester.tap(find.text('Мои данные'));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('the archive is downloaded and saved', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 1000));
    await openPrivacy(tester);

    await tester.tap(find.text('Скачать архив'));
    await tester.pumpAndSettle();

    expect(app.saver.saved.values.single, [80, 75, 3, 4]);
    expect(find.text('Архив сохранён на устройстве.'), findsOneWidget);
  });

  testWidgets('deletion is confirmed, planned and can be cancelled', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 1000));
    await openPrivacy(tester);

    await tester.tap(find.widgetWithText(TextButton, 'Удалить аккаунт'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(TextButton, 'Удалить аккаунт').last);
    await tester.pumpAndSettle();

    expect(app.server.deletion, isNotNull);
    expect(find.textContaining('Аккаунт будет удалён'), findsOneWidget);

    await tester.tap(find.text('Отменить удаление'));
    await tester.pumpAndSettle();

    expect(app.server.deletion, isNull);
  });

  testWidgets('a university that deletes accounts itself is explained', (
    tester,
  ) async {
    await pumpApp(
      tester,
      size: const Size(390, 1000),
      server: FakeServer()..selfDeletion = false,
    );
    await openPrivacy(tester);

    await tester.tap(find.widgetWithText(TextButton, 'Удалить аккаунт'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(TextButton, 'Удалить аккаунт').last);
    await tester.pumpAndSettle();

    expect(
      find.text('Удаление аккаунта в вашем вузе проводит администрация.'),
      findsOneWidget,
    );
  });
}
