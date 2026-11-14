import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_authenticator.dart';
import 'app_harness.dart';

void main() {
  testWidgets('a new install connects to the university and signs in', (
    tester,
  ) async {
    final authenticator = FakeAuthenticator();
    final app = await pumpApp(
      tester,
      size: const Size(390, 844),
      stored: const {'settings.language': 'ru', 'tour.seen': true},
      signedIn: false,
      authenticator: authenticator,
    );

    expect(find.text('Подключите свой вуз'), findsOneWidget);
    await tester.enterText(find.byType(TextField).first, 'campus.example.edu/');
    await tester.enterText(find.byType(TextField).last, 'KazNU');
    await tester.tap(find.text('Подключиться'));
    await tester.pumpAndSettle();

    expect(find.text('Вход'), findsOneWidget);
    await tester.tap(find.text('Войти через Microsoft'));
    await tester.pumpAndSettle();

    final login = authenticator.opened.single;
    expect(login.host, 'campus.example.edu');
    expect(login.path, '/v1/auth/kaznu/microsoft/login');
    expect(login.queryParameters['code_challenge_method'], 'S256');
    expect(
      app.server.requests.map((request) => request.url.path),
      contains('/v1/auth/token'),
    );
    expect(find.byType(NavigationBar), findsOneWidget);
  });

  testWidgets('an unknown university code is explained', (tester) async {
    await pumpApp(
      tester,
      size: const Size(390, 844),
      stored: const {'settings.language': 'ru', 'tour.seen': true},
      signedIn: false,
    );

    await tester.enterText(find.byType(TextField).first, 'campus.example.edu');
    await tester.enterText(find.byType(TextField).last, 'nope');
    await tester.tap(find.text('Подключиться'));
    await tester.pumpAndSettle();

    expect(
      find.text('На этом сервере нет вуза с таким кодом.'),
      findsOneWidget,
    );
  });

  testWidgets('cancelling the browser keeps the login screen', (tester) async {
    await pumpApp(
      tester,
      size: const Size(390, 844),
      signedIn: false,
      authenticator: FakeAuthenticator(code: null),
    );

    await tester.tap(find.text('Войти через Microsoft'));
    await tester.pumpAndSettle();

    expect(find.text('Вход'), findsOneWidget);
  });

  testWidgets('profile shows the account, devices and signs out', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 844));

    await tester.tap(find.byTooltip('Профиль и настройки'));
    await tester.pumpAndSettle();

    expect(find.text('Айгерим Сапарова'), findsOneWidget);
    expect(find.text('aigerim@kaznu.kz'), findsOneWidget);
    expect(find.text('Это устройство'), findsOneWidget);
    expect(find.text('Windows'), findsOneWidget);

    await tester.tap(find.byTooltip('Выйти на этом устройстве'));
    await tester.pumpAndSettle();
    expect(find.text('Windows'), findsNothing);

    await tester.scrollUntilVisible(find.text('Выйти'), 200);
    await tester.tap(find.text('Выйти'));
    await tester.pumpAndSettle();

    expect(
      app.server.requests.map((request) => request.url.path),
      contains('/v1/auth/logout'),
    );
    expect(find.text('Вход'), findsOneWidget);
  });

  testWidgets('editing the profile saves name, about and contacts', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 1200));

    await tester.tap(find.byTooltip('Профиль и настройки'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Изменить профиль'));
    await tester.pumpAndSettle();

    await tester.enterText(find.widgetWithText(TextField, 'Имя'), 'Айгерим');
    await tester.enterText(
      find.widgetWithText(TextField, 'О себе'),
      'Изучаю биологию',
    );
    await tester.tap(find.text('Добавить контакт'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.widgetWithText(TextField, 'Контакт'),
      '@aigerim',
    );
    await tester.tap(find.text('Видят все в вузе'));
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();

    expect(app.server.profile['displayName'], 'Айгерим');
    expect(app.server.profile['contacts'], [
      {'kind': 'telegram', 'value': '@aigerim', 'public': true},
    ]);
    expect(find.text('Изучаю биологию'), findsOneWidget);
  });

  testWidgets('settings stay reachable before sign-in', (tester) async {
    await pumpApp(
      tester,
      size: const Size(390, 844),
      stored: const {'settings.language': 'ru', 'tour.seen': true},
      signedIn: false,
    );

    await tester.tap(find.byTooltip('Настройки'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('English'));
    await tester.pumpAndSettle();
    await tester.pageBack();
    await tester.pumpAndSettle();

    expect(find.text('Connect your university'), findsOneWidget);
  });
}
