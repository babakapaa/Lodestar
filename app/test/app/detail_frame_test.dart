import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/app/shell/side_nav.dart';
import 'package:lodestar/features/study/schedule_screen.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';

void main() {
  testWidgets('wide detail pages keep the side panel', (tester) async {
    await pumpApp(tester, size: const Size(1440, 900));
    await tester.tap(find.text('Расписание на неделю'));
    await tester.pumpAndSettle();
    expect(find.byType(ScheduleScreen), findsOneWidget);
    expect(find.byType(SideNav), findsOneWidget);
    expect(find.byTooltip('Назад'), findsOneWidget);

    await tester.tap(find.byTooltip('Назад'));
    await tester.pumpAndSettle();
    expect(find.byType(ScheduleScreen), findsNothing);

    await tester.tap(find.text('Расписание на неделю'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Заметки').first);
    await tester.pumpAndSettle();
    expect(find.byType(ScheduleScreen), findsNothing);
  });

  testWidgets('compact detail pages stay full screen', (tester) async {
    await pumpApp(tester, size: const Size(390, 844));
    await tester.tap(find.text('Расписание на неделю'));
    await tester.pumpAndSettle();
    expect(find.byType(ScheduleScreen), findsOneWidget);
    expect(find.byType(SideNav), findsNothing);
  });
}
