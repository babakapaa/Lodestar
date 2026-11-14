import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/chats/chat_screen.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';
import 'chats_test.dart' show openChats, seeded, settleSync, write;

void main() {
  testWidgets('wide chats show the list and the open chat side by side', (
    tester,
  ) async {
    final app = await pumpApp(
      tester,
      size: const Size(1440, 900),
      server: seeded(),
    );
    await openChats(tester);
    expect(find.text('Выберите чат'), findsOneWidget);

    await tester.tap(find.text('Ерлан Ким'));
    await settleSync(tester);
    expect(find.byType(ChatScreen), findsOneWidget);
    expect(find.text('Ерлан Ким'), findsNWidgets(2));
    expect(find.text('Биология'), findsOneWidget);

    await write(tester, 'Да, вчера');
    expect(app.server.chats.sent.single['text'], 'Да, вчера');

    await tester.enterText(find.byType(TextField).first, 'Био');
    await tester.pumpAndSettle();
    expect(find.text('Биология'), findsOneWidget);
  });
}
