import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/calls/call_peer.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';
import 'chats_extras_test.dart' show grouped, openGroup;
import 'chats_test.dart' show openChats, settleSync;

void main() {
  testWidgets('a call starts from the chat and ends with hang up', (
    tester,
  ) async {
    final server = grouped();
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: server,
    );
    await openGroup(tester);

    await tester.tap(find.byKey(const ValueKey('chat-call')));
    await tester.pumpAndSettle();
    final engine = app.calls.single;
    expect(engine.url, 'wss://calls.example.edu');
    expect(engine.token, 'room-token');
    expect(engine.microphone, isTrue);
    expect(find.text('Айгерим Сапарова (вы)'), findsOneWidget);
    expect(server.chats.calls['g1']!['people'], ['u1']);

    await tester.tap(find.byKey(const ValueKey('call-camera')));
    await tester.pumpAndSettle();
    expect(engine.camera, isTrue);
    expect(find.byKey(const ValueKey('call-video-u1')), findsOneWidget);

    engine.join(
      const CallPeer(
        id: 'u2',
        name: 'Ерлан Ким',
        local: false,
        speaking: true,
        microphone: false,
        video: false,
        screen: false,
      ),
    );
    await tester.pumpAndSettle();
    expect(find.byKey(const ValueKey('call-peer-u2')), findsOneWidget);
    expect(find.byIcon(Icons.mic_off), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('call-screen')));
    await tester.pumpAndSettle();
    expect(engine.screen, isTrue);

    await tester.tap(find.byKey(const ValueKey('call-hang-up')));
    await tester.pumpAndSettle();
    expect(engine.left, isTrue);
    expect(server.chats.calls, isEmpty);
    expect(find.byKey(const ValueKey('chat-call')), findsOneWidget);
  });

  testWidgets('an ongoing call is shown in the chat and can be joined', (
    tester,
  ) async {
    final server = grouped();
    server.chats.startCall('g1', ['u2', 'u3']);
    server.chats.screenShare = false;
    final app = await pumpApp(
      tester,
      size: const Size(390, 1400),
      server: server,
    );
    await openGroup(tester);

    expect(find.byKey(const ValueKey('chat-active-call')), findsOneWidget);
    expect(find.text('2 участника'), findsOneWidget);

    await tester.tap(find.text('Войти'));
    await tester.pumpAndSettle();
    expect(app.calls.single.url, isNotNull);
    expect(server.chats.calls['g1']!['people'], ['u2', 'u3', 'u1']);
    expect(find.byKey(const ValueKey('call-screen')), findsNothing);
  });

  testWidgets('a call explains why it cannot start', (tester) async {
    final server = grouped()..chats.callError = 'calls_off';
    await pumpApp(tester, size: const Size(390, 1400), server: server);
    await openGroup(tester);

    await tester.tap(find.byKey(const ValueKey('chat-call')));
    await tester.pumpAndSettle();
    expect(find.byKey(const ValueKey('call-error')), findsOneWidget);
    expect(find.text('Вуз выключил звонки'), findsOneWidget);
  });

  testWidgets('channels have no call button', (tester) async {
    final server = FakeServer();
    server.chats.addChat(id: 'ch1', kind: 'channel', title: 'Новости');
    await pumpApp(tester, size: const Size(390, 1400), server: server);
    await openChats(tester);
    await tester.tap(find.text('Новости'));
    await settleSync(tester);
    expect(find.byKey(const ValueKey('chat-call')), findsNothing);
  });
}
