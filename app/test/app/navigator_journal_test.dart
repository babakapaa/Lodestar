import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';
import 'navigator_test.dart' show ask;
import 'search_map_test.dart' show scrollTo;
import 'study_test.dart' show goBack, openAssignment;

void main() {
  testWidgets(
    'journal queries go into the search map and a letter is drafted',
    (tester) async {
      final app = await pumpApp(tester, size: const Size(390, 1400));
      String? copied;
      tester.binding.defaultBinaryMessenger.setMockMethodCallHandler(
        SystemChannels.platform,
        (call) async {
          if (call.method == 'Clipboard.setData') {
            copied = (call.arguments as Map)['text'] as String?;
          }
          return null;
        },
      );
      await openAssignment(tester);
      await tester.tap(find.text('Спросить навигатора'));
      await tester.pumpAndSettle();
      await ask(tester, 'Что делает мембрана?');

      await tester.tap(find.byTooltip('Написать преподавателю'));
      await tester.pumpAndSettle();
      await tester.tap(find.text('Скопировать письмо'));
      await tester.pumpAndSettle();
      expect(copied, contains('«Эссе о клетке» по курсу «Биология»'));
      expect(copied, contains('Что делает мембрана?'));
      expect(copied, contains('- Иванов И. И. Биология. 2020.'));
      expect(copied, contains('Мой вопрос:\nЧто делает мембрана?'));
      await tester.tapAt(const Offset(20, 20));
      await tester.pump(const Duration(seconds: 5));
      await tester.pumpAndSettle();

      await goBack(tester);
      await tester.tap(find.text('Собрать карту поиска'));
      await tester.pumpAndSettle();
      await scrollTo(tester, find.text('Взять запросы из журнала поиска'));
      await tester.tap(find.text('Взять запросы из журнала поиска'));
      await tester.pumpAndSettle();
      expect(
        find.widgetWithText(TextField, 'Что делает мембрана?'),
        findsOneWidget,
      );
      await scrollTo(tester, find.text('Взять запросы из журнала поиска'));
      await tester.tap(find.text('Взять запросы из журнала поиска'));
      await tester.pumpAndSettle();
      expect(
        find.text('В журнале нет новых запросов по этому заданию'),
        findsOneWidget,
      );
      expect(app.server.navigator.journal, hasLength(1));
    },
  );

  testWidgets('the journal can be switched off and cleared', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 1000));
    app.server.navigator.journal.add({
      'id': 'j1',
      'scopeKind': '',
      'scopeId': '',
      'query': 'Что такое митоз?',
      'sources': <Object?>[],
      'createdAt': '2026-11-13T10:00:00Z',
    });
    await tester.tap(find.text('Навигатор').first);
    await tester.pumpAndSettle();
    await tester.tap(find.byTooltip('Журнал поиска'));
    await tester.pumpAndSettle();
    expect(find.text('Что такое митоз?'), findsOneWidget);

    await tester.tap(find.text('Вести журнал поиска'));
    await tester.pumpAndSettle();
    expect(app.server.preferences['navigator.journal'], isFalse);

    await tester.tap(find.byTooltip('Удалить запись'));
    await tester.pumpAndSettle();
    expect(app.server.navigator.journal, isEmpty);
    expect(find.text('Журнал пока пуст'), findsOneWidget);
  });

  testWidgets('the red button shows university contacts', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 1000));
    await tester.tap(find.text('Навигатор').first);
    await tester.pumpAndSettle();
    await tester.tap(find.byTooltip('Мне тяжело'));
    await tester.pumpAndSettle();
    expect(find.textContaining('Вуз пока не указал контакты'), findsOneWidget);
    await tester.tapAt(const Offset(20, 20));
    await tester.pumpAndSettle();

    app.server.navigator.support = 'Психологическая служба: 112';
    await tester.tap(find.byTooltip('Мне тяжело'));
    await tester.pumpAndSettle();
    expect(find.text('Психологическая служба: 112'), findsOneWidget);
  });
}
