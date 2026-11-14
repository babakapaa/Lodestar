import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/navigator/artifact_panel.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';
import 'navigator_test.dart' show ask;

void main() {
  testWidgets('wide navigator keeps history and opens documents beside', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(1440, 900));
    await tester.tap(find.text('Навигатор').first);
    await tester.pumpAndSettle();
    expect(find.text('О чём хотите разобраться?'), findsOneWidget);

    await ask(tester, 'Что делает мембрана?');
    expect(app.server.navigator.dialogs, hasLength(1));
    expect(find.byKey(const ValueKey('navigator-question')), findsOneWidget);

    await tester.tap(find.byTooltip('ИИ-декларация'));
    await tester.pumpAndSettle();
    expect(find.byType(ArtifactPanel), findsOneWidget);
    expect(find.textContaining('# ИИ-декларация'), findsWidgets);

    await tester.tap(find.byTooltip('Закрыть документ'));
    await tester.pumpAndSettle();
    expect(find.byType(ArtifactPanel), findsNothing);

    await tester.tap(find.text('Новый вопрос').first);
    await tester.pumpAndSettle();
    expect(find.text('О чём хотите разобраться?'), findsOneWidget);
  });
}
