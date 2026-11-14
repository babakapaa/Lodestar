import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';
import 'chats_test.dart' show seeded;

void main() {
  testWidgets('wide today shows quick actions, messages and notes', (
    tester,
  ) async {
    final server = seeded();
    server.chats.incoming('d1', 'u2', 'Привет, ты на паре?');
    await pumpApp(tester, size: const Size(1440, 900), server: server);
    await tester.pumpAndSettle();
    expect(find.text('Спросить навигатора'), findsOneWidget);
    expect(find.text('Новые сообщения'), findsOneWidget);
    expect(find.text('Недавние заметки'), findsOneWidget);
    expect(
      find.text('Заметок пока нет. Создайте первую кнопкой выше.'),
      findsOneWidget,
    );
  });
}
