import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';
import 'study_test.dart' show goBack, openAssignment;

Future<void> ask(WidgetTester tester, String text) async {
  await tester.enterText(
    find.byKey(const ValueKey('navigator-question')),
    text,
  );
  await tester.tap(find.byTooltip('Отправить'));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('a student climbs the hint ladder about an assignment', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 1400));
    await openAssignment(tester);
    await tester.tap(find.text('Спросить навигатора'));
    await tester.pumpAndSettle();
    expect(find.textContaining('не пишет работу за вас'), findsOneWidget);

    await ask(tester, 'Что делает мембрана?');
    final dialog = app.server.navigator.dialogs.values.single;
    expect(dialog['scopeKind'], 'assignment');
    expect(dialog['scopeId'], 'a1');
    expect(find.text('Задание: Эссе о клетке'), findsOneWidget);
    expect(find.text('Ступень 1 из 5: уточняющий вопрос'), findsOneWidget);
    expect(find.text('Строение клетки'), findsOneWidget);

    await tester.tap(find.text('Строение клетки'));
    await tester.pumpAndSettle();
    expect(
      find.text('Мембрана отделяет содержимое клетки от внешней среды.'),
      findsOneWidget,
    );
    await tester.tapAt(const Offset(20, 20));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Ещё подсказку'));
    await tester.pumpAndSettle();
    expect(app.server.navigator.questions.last['more'], isTrue);
    expect(find.text('Ступень 2 из 5: направление'), findsOneWidget);

    await tester.tap(find.text('Меня не так поняли'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Найти источник'));
    await tester.pumpAndSettle();
    final misread = app.server.navigator.questions.last;
    expect(misread['intent'], 'find');
    expect(misread['text'], 'Что делает мембрана?');

    await goBack(tester);
    await goBack(tester);
    await goBack(tester);
    await tester.tap(find.text('Навигатор').first);
    await tester.pumpAndSettle();
    expect(find.text('Что делает мембрана?'), findsOneWidget);

    await tester.tap(find.text('Что делает мембрана?'));
    await tester.pumpAndSettle();
    expect(find.text('Ступень 1 из 5: уточняющий вопрос'), findsNWidgets(2));
    await tester.tap(find.byTooltip('Удалить диалог'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Удалить диалог').last);
    await tester.pumpAndSettle();
    expect(app.server.navigator.dialogs, isEmpty);
    expect(find.text('Навигатор поможет разобраться в теме'), findsOneWidget);
  });

  testWidgets('a failed question is retried with the same key', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 1000));
    await tester.tap(find.text('Навигатор').first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Новый вопрос'));
    await tester.pumpAndSettle();

    app.server.navigator.unavailable = true;
    await ask(tester, 'Мне тяжело, ничего не успеваю');
    expect(
      find.text('Навигатор в вашем вузе ещё не подключён.'),
      findsOneWidget,
    );

    app.server.navigator.unavailable = false;
    await tester.tap(find.text('Повторить'));
    await tester.pumpAndSettle();
    final keys = {for (final q in app.server.navigator.questions) q['key']};
    expect(keys, hasLength(1));
    expect(find.text('Психологическая служба: 112'), findsOneWidget);
    expect(app.server.navigator.dialogs.values.single['scopeKind'], '');
  });

  testWidgets('a teacher adds course materials for the navigator', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 1400));
    app.server.study.role = 'teacher';
    await tester.tap(find.text('Учёба').first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Биология'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Материалы для навигатора'));
    await tester.pumpAndSettle();
    expect(find.text('Материалов пока нет'), findsOneWidget);

    await tester.tap(find.text('Добавить материал'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();
    expect(
      find.text('Нужны название и текст. Ссылка должна начинаться с http.'),
      findsOneWidget,
    );
    await tester.enterText(
      find.byKey(const ValueKey('material-title')),
      'Строение клетки',
    );
    await tester.enterText(
      find.byKey(const ValueKey('material-text')),
      'Мембрана отделяет клетку.',
    );
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();
    expect(app.server.navigator.materials.single['title'], 'Строение клетки');
    expect(find.textContaining('Индексируется: 0 из 2'), findsOneWidget);

    await tester.tap(find.byTooltip('Удалить материал'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Удалить материал').last);
    await tester.pumpAndSettle();
    expect(app.server.navigator.materials, isEmpty);
  });
}
