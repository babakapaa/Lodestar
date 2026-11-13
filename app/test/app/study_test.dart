import 'dart:typed_data';

import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/files/opened_file.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';

Future<void> openAssignment(WidgetTester tester) async {
  await tester.tap(find.text('Учёба').first);
  await tester.pumpAndSettle();
  await tester.tap(find.text('Биология'));
  await tester.pumpAndSettle();
  await tester.scrollUntilVisible(
    find.text('Эссе о клетке'),
    200,
    scrollable: find.byType(Scrollable).last,
  );
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

  testWidgets('a student attaches files to work and downloads them', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 1400));
    await openAssignment(tester);

    app.opener.next = [
      OpenedFile('cell.pdf', Uint8List.fromList([4, 2])),
    ];
    await tester.tap(find.byKey(const ValueKey('submit-attach')));
    await tester.pumpAndSettle();
    expect(find.text('cell.pdf'), findsOneWidget);

    await tester.tap(find.text('Сдать работу'));
    await tester.pumpAndSettle();
    final submission = app.server.study.submissions.single;
    expect(submission['body'], '');
    expect(((submission['files']! as List).single as Map)['name'], 'cell.pdf');

    await tester.tap(find.text('cell.pdf'));
    await tester.pumpAndSettle();
    expect(app.saver.saved['cell.pdf'], Uint8List.fromList([4, 2]));
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
    await pumpApp(tester, size: const Size(390, 1700), server: server);
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
    await tester.tap(find.text('Короткая беседа'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Сохранить'));
    await tester.pumpAndSettle();

    final created = server.study.assignments.last;
    expect(created['defense'], 'short');
    expect(created['title'], 'Лабораторная 1');
    expect(created['grading'], 'pass');
    expect(created['maxPoints'], isNull);
    expect(created['publishedAt'], isNotNull);

    await tester.tap(find.text('Лабораторная 1'));
    await tester.pumpAndSettle();
    expect(find.text('Ерлан Ким'), findsOneWidget);
    expect(find.text('Не сдано'), findsOneWidget);
  });

  testWidgets('a teacher builds the course schedule and changes a class', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    await pumpApp(tester, size: const Size(390, 1700), server: server);
    await tester.tap(find.text('Учёба').first);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Биология'));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('course-schedule')));
    await tester.pumpAndSettle();
    expect(find.text('Занятий пока нет. Добавьте первое.'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('series-add')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('series-time')), '9:00');
    await tester.tap(find.byKey(const ValueKey('series-save')));
    await tester.pumpAndSettle();
    expect(find.text('Проверьте время, длительность и даты'), findsOneWidget);
    await tester.enterText(find.byKey(const ValueKey('series-time')), '09:00');
    await tester.enterText(find.byKey(const ValueKey('series-room')), 'А-204');
    await tester.tap(find.byKey(const ValueKey('series-save')));
    await tester.pumpAndSettle();

    final created = server.study.series.single;
    expect(created['startTime'], '09:00');
    expect(created['room'], 'А-204');
    expect(created['minutes'], 90);
    expect(find.byKey(ValueKey('series-${created['id']}')), findsOneWidget);

    final lessonKey = ValueKey('lesson-${created['id']}-2026-11-16');
    await tester.tap(find.byKey(lessonKey));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('lesson-cancel')));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byKey(const ValueKey('team-text')),
      'Конференция',
    );
    await tester.tap(find.text('Отменить занятие').last);
    await tester.pumpAndSettle();
    expect(server.study.courseLessons.single['status'], 'cancelled');
    expect(find.textContaining('Конференция'), findsOneWidget);

    await tester.tap(find.byKey(lessonKey));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('lesson-revert')));
    await tester.pumpAndSettle();
    expect(server.study.courseLessons.single['status'], 'planned');
  });

  testWidgets('a teacher grades submitted work', (tester) async {
    final server = FakeServer();
    server.study.role = 'teacher';
    server.study.submissions.add({
      'id': 's1',
      'assignmentId': 'a1',
      'clientId': 'x',
      'body': 'Ядро, мембрана, цитоплазма',
      'submittedAt': '2026-11-11T10:00:00Z',
      'files': <Object?>[],
    });
    await pumpApp(tester, size: const Size(390, 1400), server: server);
    await openAssignment(tester);

    await tester.tap(find.text('Ерлан Ким'));
    await tester.pumpAndSettle();
    expect(find.text('Ядро, мембрана, цитоплазма'), findsOneWidget);

    await tester.enterText(
      find.widgetWithText(TextField, 'Баллы, максимум 10'),
      '11',
    );
    await tester.tap(find.text('Сохранить оценку'));
    await tester.pumpAndSettle();
    expect(
      find.text('Проверьте оценку. Без оценки нужен комментарий.'),
      findsOneWidget,
    );

    await tester.enterText(
      find.widgetWithText(TextField, 'Баллы, максимум 10'),
      '8,5',
    );
    await tester.enterText(
      find.widgetWithText(TextField, 'Комментарий'),
      'Хорошо, но мало про митохондрии',
    );
    await tester.tap(find.text('Сохранить оценку'));
    await tester.pumpAndSettle();

    expect(server.study.grades['a1']!['points'], 8.5);
    expect(server.study.grades['a1']!['passed'], isNull);
    expect(find.text('8,5 из 10'), findsOneWidget);
  });
}
