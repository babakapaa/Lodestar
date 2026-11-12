import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';

void main() {
  testWidgets('phone shows bottom navigation with five tabs', (tester) async {
    await pumpApp(tester, size: const Size(390, 844));

    expect(find.byType(NavigationBar), findsOneWidget);
    expect(find.byType(NavigationRail), findsNothing);
    expect(find.byType(NavigationDestination), findsNWidgets(5));
    expect(find.text('Сегодня свободно'), findsOneWidget);
  });

  testWidgets('tablet shows a compact side rail', (tester) async {
    await pumpApp(tester, size: const Size(800, 1000));

    final rail = tester.widget<NavigationRail>(find.byType(NavigationRail));
    expect(rail.extended, isFalse);
    expect(find.byType(NavigationBar), findsNothing);
  });

  testWidgets('desktop shows the full side panel with the name', (
    tester,
  ) async {
    await pumpApp(tester, size: const Size(1280, 800));

    final rail = tester.widget<NavigationRail>(find.byType(NavigationRail));
    expect(rail.extended, isTrue);
    expect(find.text('Lodestar'), findsOneWidget);
    expect(find.text('Профиль и настройки'), findsOneWidget);
  });

  testWidgets('tapping a tab opens its screen', (tester) async {
    await pumpApp(tester, size: const Size(390, 844));

    await tester.tap(find.text('Чаты'));
    await tester.pumpAndSettle();

    expect(find.text('Чаты курса и команд'), findsOneWidget);
  });

  testWidgets('Kazakh is used when chosen', (tester) async {
    await pumpApp(
      tester,
      size: const Size(390, 844),
      stored: const {
        'settings.language': 'kk',
        'connection.server': 'https://campus.example.edu',
        'connection.tenant': 'kaznu',
      },
    );

    expect(find.text('Бүгін'), findsWidgets);
  });
}
