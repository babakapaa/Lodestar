import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'admin_navigation.dart';
import 'app_harness.dart';

void main() {
  testWidgets('the rector sees units without small groups', (tester) async {
    final server = FakeServer()..permissions = ['core.stats.view'];
    await pumpApp(tester, size: const Size(390, 1200), server: server);
    await openAdmin(tester, 'Статистика вуза');

    expect(find.text('Биофак'), findsOneWidget);
    expect(find.text('Людей: 120'), findsOneWidget);
    expect(find.text('82%'), findsOneWidget);
    expect(find.text('12.5 ч'), findsOneWidget);
    expect(find.text('Меньше 5 человек, данные скрыты'), findsOneWidget);
    expect(find.textContaining('Данных об отдельных людях'), findsOneWidget);

    await tester.tap(find.text('Год'));
    await tester.pumpAndSettle();
    expect(server.stats.queries.last['days'], '365');

    await tester.tap(find.text('Худграф'));
    await tester.pumpAndSettle();
    expect(find.text('Кафедра генетики'), findsNothing);

    await tester.tap(find.text('Биофак'));
    await tester.pumpAndSettle();
    expect(server.stats.queries.last['parent'], 'bio');
    expect(find.text('Кафедра генетики'), findsOneWidget);
    expect(find.text('Нет данных за этот период'), findsOneWidget);
  });

  testWidgets('university settings are grouped by section', (tester) async {
    final server = FakeServer()..permissions = ['core.administrator'];
    await pumpApp(tester, size: const Size(390, 1200), server: server);
    await openAdmin(tester, 'Настройки вуза');

    expect(find.text('Общее'), findsOneWidget);
    expect(find.text('Право и данные'), findsOneWidget);
    expect(find.text('Песочница GPU'), findsNothing);
  });
}
