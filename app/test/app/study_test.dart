import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';

Future<void> openAssignment(WidgetTester tester) async {
  await tester.tap(find.text('Учёба').first);
  await tester.pumpAndSettle();
  await tester.tap(find.text('Биология'));
  await tester.pumpAndSettle();
  await tester.tap(find.text('Эссе о клетке'));
  await tester.pumpAndSettle();
}

Future<void> goBack(WidgetTester tester) async {
  await tester.tap(find.byTooltip('Назад'));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('a student submits work online', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 1400));
    await openAssignment(tester);

    expect(find.text('Опишите строение клетки.'), findsOneWidget);
    expect(find.text('Вы ещё не сдавали эту работу.'), findsOneWidget);

    await tester.enterText(find.byType(TextField), 'Клетка состоит из ядра');
    await tester.tap(find.text('Сдать работу'));
    await tester.pumpAndSettle();

    final submission = app.server.study.submissions.single;
    expect(submission['body'], 'Клетка состоит из ядра');
    expect(submission['clientId'], isNotEmpty);
    expect(find.text('Клетка состоит из ядра'), findsOneWidget);
    expect(find.text('Сдать заново'), findsOneWidget);
  });

  testWidgets('work written offline is sent when the network returns', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 1400));
    await openAssignment(tester);

    app.server.offline = true;
    await tester.enterText(find.byType(TextField), 'Ответ без сети');
    await tester.tap(find.text('Сдать работу'));
    await tester.pumpAndSettle();

    expect(app.server.study.submissions, isEmpty);
    expect(
      find.text('Работа сохранена и уйдёт, когда появится сеть.'),
      findsOneWidget,
    );

    await goBack(tester);
    expect(find.textContaining('Ждёт отправки'), findsOneWidget);
    await goBack(tester);
    await tester.fling(find.text('Биология'), const Offset(0, 400), 1000);
    await tester.pumpAndSettle();
    expect(find.text('Нет сети. Показаны сохранённые данные.'), findsOneWidget);
    expect(find.text('Отправится, когда появится сеть'), findsOneWidget);

    app.server.offline = false;
    app.live.send(1, 'core.ping', null);
    await tester.pumpAndSettle();

    expect(app.server.study.submissions.single['body'], 'Ответ без сети');
    expect(find.text('Ждут отправки'), findsNothing);
  });

  testWidgets('a draft retries on a timer and keeps one copy', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 1400));
    await openAssignment(tester);

    app.server.offline = true;
    await tester.enterText(find.byType(TextField), 'Повторная отправка');
    await tester.tap(find.text('Сдать работу'));
    await tester.pumpAndSettle();
    app.server.offline = false;

    await tester.pump(const Duration(seconds: 6));
    await tester.pumpAndSettle();
    expect(app.server.study.submissions, hasLength(1));
    expect(
      find.text('Работа сохранена и уйдёт, когда появится сеть.'),
      findsNothing,
    );
  });

  testWidgets('a rejected draft can be discarded', (tester) async {
    final server = FakeServer();
    server.study.submitError = 'assignment_not_published';
    await pumpApp(tester, size: const Size(390, 1400), server: server);
    await openAssignment(tester);

    await tester.enterText(find.byType(TextField), 'Поздно');
    await tester.tap(find.text('Сдать работу'));
    await tester.pumpAndSettle();
    await goBack(tester);
    await goBack(tester);

    expect(
      find.text('Не отправилось: Задание ещё не опубликовано.'),
      findsOneWidget,
    );
    await tester.tap(find.byTooltip('Удалить'));
    await tester.pumpAndSettle();
    expect(find.text('Ждут отправки'), findsNothing);
  });

  testWidgets('a teacher publishes an assignment and sees student work', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    await pumpApp(tester, size: const Size(390, 1400), server: server);
    await tester.tap(find.text('Учёба').first);
    await tester.pumpAndSettle();
    expect(find.text('BIO-101, Преподаю'), findsOneWidget);
    await tester.tap(find.text('Биология'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Новое задание'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.widgetWithText(TextField, 'Название'),
      'Лабораторная 1',
    );
    await tester.tap(find.text('Зачёт'));
    await tester.tap(find.text('Опубликовать'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();

    final created = server.study.assignments.last;
    expect(created['title'], 'Лабораторная 1');
    expect(created['grading'], 'pass');
    expect(created['maxPoints'], isNull);
    expect(created['publishedAt'], isNotNull);

    await tester.tap(find.text('Лабораторная 1'));
    await tester.pumpAndSettle();
    expect(find.text('Ерлан Ким'), findsOneWidget);
    expect(find.text('Не сдано'), findsOneWidget);
  });
}
