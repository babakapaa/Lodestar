import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';

void main() {
  testWidgets('choices made on another device arrive after sign-in', (
    tester,
  ) async {
    final server = FakeServer()..preferences['core.language'] = 'en';

    await pumpApp(tester, size: const Size(390, 844), server: server);

    expect(find.text('Today'), findsWidgets);
  });

  testWidgets('a change here is sent to the server', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 844));

    await tester.tap(find.byTooltip('Профиль и настройки'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Настройки'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Тёмная'));
    await tester.pumpAndSettle();

    expect(app.server.preferences['core.appearance'], 'dark');
  });
}
