import 'dart:typed_data';

import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/files/opened_file.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'admin_navigation.dart';
import 'app_harness.dart';
import 'chats_test.dart' show openChats, settleSync;
import 'study_test.dart' show goBack;

FakeServer grouped() {
  final server = FakeServer();
  server.chats.addChat(
    id: 'g1',
    kind: 'group',
    title: 'Проект ИС-21',
    role: 'owner',
    people: {'u2': 'Ерлан Ким', 'u3': 'Дана Ахметова'},
  );
  server.files.add('f9', 'plan.pdf', [1, 2, 3]);
  server.chats.incoming('g1', 'u3', 'План проекта', fileIds: ['f9']);
  return server;
}

Future<void> openGroup(WidgetTester tester) async {
  await openChats(tester);
  await tester.tap(find.text('Проект ИС-21'));
  await settleSync(tester);
}

void main() {
  testWidgets('files are attached, downloaded and people mentioned', (
    tester,
  ) async {
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: grouped(),
    );
    await openGroup(tester);

    await tester.tap(find.text('plan.pdf'));
    await tester.pumpAndSettle();
    expect(app.saver.saved['plan.pdf'], Uint8List.fromList([1, 2, 3]));

    app.opener.next = [
      OpenedFile('essay.pdf', Uint8List.fromList([7, 7])),
    ];
    await tester.tap(find.byKey(const ValueKey('chat-attach')));
    await tester.pumpAndSettle();
    expect(find.text('essay.pdf'), findsOneWidget);

    final field = find.byKey(const ValueKey('chat-message'));
    await tester.enterText(field, '@Ер');
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('mention-u2')));
    await tester.pumpAndSettle();
    expect(find.text('@Ерлан Ким '), findsOneWidget);
    await tester.enterText(field, '@Ерлан Ким посмотри эссе');
    await tester.tap(find.byTooltip('Отправить'));
    await tester.pump();
    await settleSync(tester);

    final sent = app.server.chats.sent.single;
    expect(sent['mentions'], ['u2']);
    expect(sent['fileIds'], ['f1']);
    expect(app.server.files.stored['f1'], Uint8List.fromList([7, 7]));
    expect(find.text('essay.pdf'), findsOneWidget);
  });

  testWidgets('a report carries only the chosen message to moderators', (
    tester,
  ) async {
    final server = grouped()
      ..chats.incoming('g1', 'u2', 'Скинь ответы на экзамен');
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: server,
    );
    await openGroup(tester);

    await tester.longPress(find.text('Скинь ответы на экзамен'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Пожаловаться'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('chat-report-reason')),
      'Списывание',
    );
    await tester.tap(find.text('Отправить жалобу'));
    await tester.pumpAndSettle();
    expect(find.text('Жалоба отправлена'), findsOneWidget);
    final report = app.server.chats.reports.single;
    expect(report['reason'], 'Списывание');
    expect(report['messages'], hasLength(1));

    await goBack(tester);
    app.server.permissions = ['chat.reports.review'];
    await openAdmin(tester, 'Жалобы в чатах');
    expect(find.text('Причина: Списывание'), findsOneWidget);
    expect(find.text('Скинь ответы на экзамен'), findsOneWidget);
    expect(find.text('План проекта'), findsNothing);
    await tester.tap(find.text('Разобрано'));
    await tester.pumpAndSettle();
    expect(find.text('Жалоб нет'), findsOneWidget);
  });

  testWidgets('owners manage the chat, members leave, peers get blocked', (
    tester,
  ) async {
    final server = grouped();
    server.chats.addChat(
      id: 'g2',
      kind: 'group',
      title: 'Староста',
      people: {'u3': 'Дана Ахметова'},
    );
    server.chats.addChat(
      id: 'd1',
      kind: 'direct',
      title: 'Ерлан Ким',
      people: {'u2': 'Ерлан Ким'},
    );
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: server,
    );
    await openGroup(tester);
    await tester.tap(find.byKey(const ValueKey('chat-info')));
    await tester.pumpAndSettle();

    await tester.tap(find.byKey(const ValueKey('chat-edit-info')));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('chat-edit-about')),
      'Курсовой проект',
    );
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();
    expect(app.server.chats.chats['g1']!['about'], 'Курсовой проект');
    expect(find.text('Курсовой проект'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('remove-u3')));
    await tester.pumpAndSettle();
    expect(app.server.chats.members['g1']!.containsKey('u3'), isFalse);
    expect(find.text('Дана Ахметова'), findsNothing);

    await tester.tap(find.byKey(const ValueKey('chat-mute')));
    await tester.pumpAndSettle();
    expect(app.server.chats.chats['g1']!['muted'], isTrue);
    expect(find.byKey(const ValueKey('chat-leave')), findsNothing);

    await goBack(tester);
    await goBack(tester);
    await tester.tap(find.text('Староста'));
    await settleSync(tester);
    await tester.tap(find.byKey(const ValueKey('chat-info')));
    await tester.pumpAndSettle();
    expect(find.byKey(const ValueKey('remove-u3')), findsNothing);
    await tester.tap(find.byKey(const ValueKey('chat-leave')));
    await settleSync(tester);
    expect(app.server.chats.chats.containsKey('g2'), isFalse);
    expect(find.text('Староста'), findsNothing);

    await tester.tap(find.text('Ерлан Ким'));
    await settleSync(tester);
    await tester.tap(find.byKey(const ValueKey('chat-info')));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Заблокировать'));
    await tester.pumpAndSettle();
    expect(app.server.chats.blocked, {'u2'});
    expect(find.text('Разблокировать'), findsOneWidget);
  });

  testWidgets('course chats come from the timetable, teachers add a group', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    server.chats.addChat(
      id: 'cc1',
      kind: 'channel',
      title: 'Биология',
      origin: {'kind': 'study.course', 'id': 'c1'},
    );
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: server,
    );
    await settleSync(tester);
    await tester.tap(find.text('Учёба').first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Биология'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Чаты курса'));
    await tester.pumpAndSettle();
    expect(find.text('Биология'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('course-group')));
    await settleSync(tester);
    final created = app.server.chats.chats.values.last;
    expect(created['kind'], 'group');
    expect(created['origin'], {'kind': 'study.course', 'id': 'c1'});
  });
}
