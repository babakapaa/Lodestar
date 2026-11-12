import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/design/tokens/family_specs.dart';
import 'package:lodestar/design/tokens/lodestar_tokens.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';

LodestarTokens tokensOf(WidgetTester tester) {
  final context = tester.element(find.byType(Scaffold).last);
  return Theme.of(context).extension<LodestarTokens>()!;
}

Future<void> openSettings(WidgetTester tester) async {
  await tester.tap(find.byTooltip('Профиль и настройки'));
  await tester.pumpAndSettle();
  await tester.tap(find.text('Настройки'));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('choosing a theme and dark mode restyles the app', (
    tester,
  ) async {
    await pumpApp(tester, size: const Size(390, 844));

    await openSettings(tester);
    expect(find.text('Настройки'), findsOneWidget);

    await tester.tap(find.text('Лаборатория'));
    await tester.pumpAndSettle();
    expect(tokensOf(tester).palette.accent, FamilySpecs.lab.light.accent);

    await tester.tap(find.text('Тёмная'));
    await tester.pumpAndSettle();
    expect(tokensOf(tester).palette.bg, FamilySpecs.lab.dark.bg);
  });

  testWidgets('switching language changes the interface at once', (
    tester,
  ) async {
    await pumpApp(tester, size: const Size(390, 844));

    await openSettings(tester);
    await tester.tap(find.text('English'));
    await tester.pumpAndSettle();

    expect(find.text('Settings'), findsOneWidget);
  });
}
