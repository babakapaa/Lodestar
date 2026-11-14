import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/tour/tour_dialog.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';

void main() {
  testWidgets('the tour opens on first run and can be replayed', (
    tester,
  ) async {
    final stored = Map<String, Object>.of(signedInStorage)..remove('tour.seen');
    await pumpApp(tester, size: const Size(1440, 900), stored: stored);
    expect(find.byType(TourDialog), findsOneWidget);
    expect(find.text('«Сегодня» — главный экран'), findsOneWidget);

    await tester.tap(find.text('Далее'));
    await tester.pumpAndSettle();
    expect(
      find.text('Навигатор подсказывает, а не пишет за вас'),
      findsOneWidget,
    );

    await tester.tap(find.text('Назад'));
    await tester.pumpAndSettle();
    expect(find.text('«Сегодня» — главный экран'), findsOneWidget);

    for (var step = 0; step < 4; step++) {
      await tester.tap(find.text('Далее'));
      await tester.pumpAndSettle();
    }
    expect(find.text('Пропустить'), findsNothing);
    await tester.tap(find.text('Начать'));
    await tester.pumpAndSettle();
    expect(find.byType(TourDialog), findsNothing);

    await tester.tap(find.text('Профиль и настройки').first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Настройки').first);
    await tester.pumpAndSettle();
    await tester.scrollUntilVisible(find.text('Показать тур снова'), 200);
    await tester.tap(find.text('Показать тур снова'));
    await tester.pumpAndSettle();
    expect(find.byType(TourDialog), findsOneWidget);

    await tester.tap(find.text('Пропустить'));
    await tester.pumpAndSettle();
    expect(find.byType(TourDialog), findsNothing);
  });

  testWidgets('a seen tour stays closed', (tester) async {
    await pumpApp(tester, size: const Size(390, 844));
    expect(find.byType(TourDialog), findsNothing);
  });
}
