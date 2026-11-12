import 'package:flutter_test/flutter_test.dart';

Future<void> openAdmin(WidgetTester tester, String section) async {
  await tester.tap(find.byTooltip('Профиль и настройки'));
  await tester.pumpAndSettle();
  await tester.tap(find.text('Администрирование'));
  await tester.pumpAndSettle();
  await tester.tap(find.text(section));
  await tester.pumpAndSettle();
}
