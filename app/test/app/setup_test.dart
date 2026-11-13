import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'admin_navigation.dart';
import 'app_harness.dart';
import 'study_test.dart' show goBack;

void main() {
  testWidgets('the owner goes through the wizard and launches', (tester) async {
    final server = FakeServer()..permissions = ['core.administrator'];
    server.setup.launched = false;
    final app = await pumpApp(
      tester,
      size: const Size(390, 1600),
      server: server,
    );

    expect(find.textContaining('Вуз ещё не запущен'), findsOneWidget);
    await openAdmin(tester, 'Мастер установки');

    expect(find.text('Утверждено 0 из 7'), findsOneWidget);
    expect(find.text('Протокол подключения не подписан'), findsOneWidget);
    await tester.tap(find.byKey(const ValueKey('setup-launch-action')));
    await tester.pumpAndSettle();
    expect(server.setup.launched, isFalse);

    await tester.tap(find.text('Право и безопасность').first);
    await tester.pumpAndSettle();
    for (final q in ['storage', 'consent']) {
      await tester.enterText(
        find.byKey(ValueKey('legal-answer-$q')),
        'Серверы вуза',
      );
      await tester.enterText(
        find.byKey(ValueKey('legal-basis-$q')),
        'Приказ ректора №1',
      );
      await tester.tap(find.byKey(ValueKey('legal-save-$q')));
      await tester.pumpAndSettle();
      await tester.tap(find.byKey(ValueKey('legal-approve-$q')));
      await tester.pumpAndSettle();
    }
    expect(server.setup.answers['storage']?['basis'], 'Приказ ректора №1');
    expect(find.text('Утверждено'), findsNWidgets(2));
    await tester.tap(find.byKey(const ValueKey('setup-sign')));
    await tester.pumpAndSettle();
    expect(find.text('Версия 1, подписан'), findsOneWidget);
    await tester.tap(find.byKey(const ValueKey('setup-protocol-pdf')));
    await tester.pumpAndSettle();
    expect(app.saver.saved.keys, contains('protocol.pdf'));
    await tester.ensureVisible(find.byKey(const ValueKey('setup-approve')));
    await tester.tap(find.byKey(const ValueKey('setup-approve')));
    await tester.pumpAndSettle();
    expect(server.setup.statuses['legal'], 'approved');

    for (final key in ['deploy', 'university', 'people', 'admins']) {
      server.setup.statuses[key] = 'approved';
    }
    await goBack(tester);
    await tester.drag(find.byType(ListView).last, const Offset(0, 400));
    await tester.pumpAndSettle();
    expect(find.text('Всё готово к запуску'), findsOneWidget);
    await tester.ensureVisible(
      find.byKey(const ValueKey('setup-launch-action')),
    );
    await tester.tap(find.byKey(const ValueKey('setup-launch-action')));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('launch-pilot')));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Основы биологии'));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('launch-confirm')));
    await tester.pumpAndSettle();
    expect(server.setup.launched, isTrue);
    expect(server.setup.pilot, ['c-bio']);
    expect(find.text('Вуз запущен'), findsOneWidget);
    expect(find.text('Открыт на 1 курсе'), findsOneWidget);
    await tester.tap(find.byKey(const ValueKey('setup-expand')));
    await tester.pumpAndSettle();
    expect(server.setup.pilot, isEmpty);
    expect(find.byKey(const ValueKey('setup-expand')), findsNothing);
  });

  testWidgets('a section setting is saved through the wizard', (tester) async {
    final server = FakeServer()..permissions = ['core.administrator'];
    server.setup.launched = false;
    await pumpApp(tester, size: const Size(390, 1600), server: server);
    await openAdmin(tester, 'Мастер установки');

    await tester.tap(find.text('Данные и приватность'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Люди сами удаляют свой аккаунт'));
    await tester.pumpAndSettle();
    expect(
      server.setup.calls,
      contains('PUT sections/privacy/settings/core.self_deletion'),
    );
    await tester.tap(find.byKey(const ValueKey('setup-submit')));
    await tester.pumpAndSettle();
    expect(server.setup.statuses['privacy'], 'review');
  });

  testWidgets('a student sees that the university is being set up', (
    tester,
  ) async {
    final server = FakeServer();
    server.setup
      ..launched = false
      ..manages = false;
    await pumpApp(tester, size: const Size(390, 1000), server: server);

    expect(find.textContaining('Вуз ещё не запущен'), findsNothing);
  });

  testWidgets('a student outside the pilot sees why the app is closed', (
    tester,
  ) async {
    final server = FakeServer();
    server.setup
      ..manages = false
      ..pilot = ['c-chem']
      ..included = false;
    await pumpApp(tester, size: const Size(390, 1000), server: server);

    expect(find.byKey(const ValueKey('setup-pilot-banner')), findsOneWidget);
  });
}
