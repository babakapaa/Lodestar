import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';

void main() {
  testWidgets('the bell shows unread and opens rendered notifications', (
    tester,
  ) async {
    await pumpApp(tester, size: const Size(390, 844));

    expect(find.text('1'), findsOneWidget);
    await tester.tap(find.byTooltip('Уведомления, непрочитанных: 1'));
    await tester.pumpAndSettle();

    expect(find.text('Скоро срок: Биология'), findsOneWidget);

    await tester.tap(find.byTooltip('Прочитать все'));
    await tester.pumpAndSettle();

    expect(find.byTooltip('Прочитать все'), findsNothing);
  });

  testWidgets('a live notification appears at once', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 844));
    await tester.tap(find.byTooltip('Уведомления, непрочитанных: 1'));
    await tester.pumpAndSettle();

    app.live.send(7, 'core.notification', {
      'id': 'n2',
      'type': 'study.deadline',
      'params': {'course': 'Химия'},
      'link': '',
      'createdAt': '2026-11-12T11:00:00Z',
      'readAt': null,
    });
    await tester.pumpAndSettle();

    expect(find.text('Скоро срок: Химия'), findsOneWidget);
    final uri = app.live.opened.single;
    expect(uri.scheme, 'wss');
    expect(uri.path, '/v1/live');
    expect(uri.queryParameters['access_token'], 'access-1');
  });

  testWidgets('the live stream reconnects after the last event', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 844));

    app.live.send(41, 'core.ping', null);
    await tester.pump();
    await app.live.drop();
    await tester.pump(const Duration(seconds: 2));

    expect(app.live.opened, hasLength(2));
    expect(app.live.opened.last.queryParameters['after'], '41');
  });

  testWidgets('channels and quiet hours are saved on the server', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 1200));
    await tester.tap(find.byTooltip('Уведомления, непрочитанных: 1'));
    await tester.pumpAndSettle();
    await tester.tap(find.byTooltip('Настройки уведомлений'));
    await tester.pumpAndSettle();

    expect(find.text('Сроки сдачи'), findsOneWidget);
    await tester.tap(find.text('На почту'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Не беспокоить ночью'));
    await tester.pumpAndSettle();

    expect(app.server.channels['study.deadline']!['email'], isTrue);
    expect(app.server.preferences['core.quiet_hours'], isTrue);
  });
}
