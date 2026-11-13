import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';
import 'study_test.dart' show goBack;

Future<void> openChats(WidgetTester tester) async {
  await tester.tap(find.text('Чаты').first);
  await tester.pumpAndSettle();
  await settleSync(tester);
}

Future<void> settleSync(WidgetTester tester) async {
  await tester.pump(const Duration(milliseconds: 200));
  await tester.pumpAndSettle();
}

Future<void> write(WidgetTester tester, String text) async {
  await tester.enterText(find.byKey(const ValueKey('chat-message')), text);
  await tester.tap(find.byTooltip('Отправить'));
  await tester.pump();
  await settleSync(tester);
}

FakeServer seeded() {
  final server = FakeServer();
  server.chats.addChat(
    id: 'd1',
    kind: 'direct',
    title: 'Ерлан Ким',
    people: {'u2': 'Ерлан Ким'},
  );
  server.chats.incoming('d1', 'u2', 'Привет! Ты сдал эссе?');
  server.chats.addChat(
    id: 'ch1',
    kind: 'channel',
    title: 'Биология',
    people: {'u3': 'Дана Ахметова'},
  );
  server.chats.incoming('ch1', 'u3', 'Семинар переносится на пятницу');
  return server;
}

void main() {
  testWidgets('a student reads, writes, edits and gets replies live', (
    tester,
  ) async {
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: seeded(),
    );
    await openChats(tester);
    expect(find.text('Ерлан Ким'), findsOneWidget);
    expect(find.text('Привет! Ты сдал эссе?'), findsOneWidget);
    expect(find.text('Семинар переносится на пятницу'), findsOneWidget);

    await tester.tap(find.text('Ерлан Ким'));
    await settleSync(tester);
    expect(find.text('Привет! Ты сдал эссе?'), findsOneWidget);
    expect(app.server.chats.chats['d1']!['unread'], 0);

    await write(tester, 'Да, вчера');
    expect(app.server.chats.sent.single['text'], 'Да, вчера');
    expect(find.text('Да, вчера'), findsOneWidget);

    app.server.chats.incoming('d1', 'u2', 'Отлично!');
    app.live.send(1, 'chat.message', {'chatId': 'd1', 'seq': 3, 'rev': 3});
    await settleSync(tester);
    expect(find.text('Отлично!'), findsOneWidget);

    await tester.longPress(find.text('Да, вчера'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Изменить'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('chat-edit')),
      'Да, ещё вчера',
    );
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();
    expect(find.text('Да, ещё вчера'), findsOneWidget);
    expect(find.text('изменено'), findsOneWidget);

    await tester.longPress(find.text('Отлично!'));
    await tester.pumpAndSettle();
    expect(find.text('Изменить'), findsNothing);
    await tester.tap(find.text('Ответить'));
    await tester.pumpAndSettle();
    await write(tester, 'Спасибо');
    expect(app.server.chats.sent.last['replyTo'], 3);
  });

  testWidgets('messages written offline go out when the network is back', (
    tester,
  ) async {
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: seeded(),
    );
    await openChats(tester);
    await tester.tap(find.text('Ерлан Ким'));
    await settleSync(tester);

    app.server.offline = true;
    await write(tester, 'Ты тут?');
    expect(find.text('Ты тут?'), findsOneWidget);
    expect(find.byIcon(Icons.schedule), findsOneWidget);
    expect(app.server.chats.sent, isEmpty);

    await goBack(tester);
    expect(
      find.text(
        'Нет сети. Показаны сохранённые сообщения, новые уйдут при подключении.',
      ),
      findsOneWidget,
    );
    await tester.tap(find.text('Ерлан Ким'));
    await tester.pumpAndSettle();
    expect(find.text('Привет! Ты сдал эссе?'), findsOneWidget);

    app.server.offline = false;
    await tester.pump(const Duration(seconds: 6));
    await tester.pumpAndSettle();
    expect(app.server.chats.sent.single['text'], 'Ты тут?');
    expect(find.byIcon(Icons.schedule), findsNothing);
  });

  testWidgets('channels are read-only, groups are created, requests accepted', (
    tester,
  ) async {
    final server = seeded();
    server.chats.addChat(
      id: 'r1',
      kind: 'direct',
      title: 'Незнакомец',
      request: true,
      people: {'u9': 'Незнакомец'},
    );
    server.chats.incoming('r1', 'u9', 'Продашь конспект?');
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: server,
    );
    await openChats(tester);
    expect(find.text('Незнакомец'), findsNothing);

    await tester.tap(find.text('Биология'));
    await settleSync(tester);
    expect(find.text('В канале пишут только администраторы'), findsOneWidget);
    expect(find.byKey(const ValueKey('chat-message')), findsNothing);
    await goBack(tester);

    await tester.tap(find.text('Запросы 1'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Незнакомец'));
    await settleSync(tester);
    await tester.tap(find.text('Принять'));
    await tester.pumpAndSettle();
    expect(app.server.chats.chats['r1']!['request'], false);
    expect(find.text('Принять'), findsNothing);
    await goBack(tester);

    await tester.tap(find.text('Новый чат'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Новая группа'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('chat-title')),
      'Проект ИС-21',
    );
    await tester.tap(find.text('Создать'));
    await settleSync(tester);
    expect(find.text('Проект ИС-21'), findsOneWidget);
    await write(tester, 'Всем привет');
    expect(app.server.chats.sent.single['text'], 'Всем привет');
  });
}
