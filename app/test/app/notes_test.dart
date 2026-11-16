import 'dart:convert';

import 'package:archive/archive.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/files/opened_file.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';

Future<void> openNotes(WidgetTester tester) async {
  await tester.tap(find.text('Заметки').first);
  await tester.pumpAndSettle();
}

Future<void> settle(WidgetTester tester) async {
  await tester.pump(const Duration(seconds: 2));
  await tester.pumpAndSettle();
}

Future<void> goBack(WidgetTester tester) async {
  await tester.tap(find.byTooltip('Назад'));
  await tester.pumpAndSettle();
}

Future<void> newNote(WidgetTester tester, String text) async {
  await tester.tap(find.text('Новая заметка'));
  await tester.pumpAndSettle();
  await tester.enterText(find.byType(TextField), text);
  await settle(tester);
}

Map<String, Object?> sharedNote(String role) => {
  'id': 'n-shared',
  'ownerId': 'u2',
  'role': role,
  'meta': {
    'path': 'Ерлан/Общая',
    'tags': <String>[],
    'links': <String>[],
    'assignmentId': null,
    'lessonId': null,
  },
  'seq': 0,
  'bytes': 0,
  'createdAt': '2026-11-12T10:00:00Z',
  'changedAt': '2026-11-12T10:00:00Z',
  'deletedAt': null,
};

void main() {
  testWidgets('a new note is written, renamed and synced', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 900));
    await openNotes(tester);
    expect(find.text('Заметки и карта поиска'), findsOneWidget);

    await newNote(tester, 'Клетка состоит из ядра.');
    await tester.tap(find.byType(PopupMenuButton<String>));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Переместить в папку'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.descendant(
        of: find.byType(AlertDialog),
        matching: find.byType(TextField),
      ),
      'Биология / Клетка',
    );
    await tester.tap(find.text('Сохранить'));
    await settle(tester);

    expect(find.text('Клетка'), findsOneWidget);
    final note = app.server.notes.notes.values.single;
    expect((note['meta']! as Map<String, Object?>)['path'], 'Биология/Клетка');
    expect(app.server.notes.updates[note['id']], isNotEmpty);

    await goBack(tester);
    expect(find.text('Клетка'), findsOneWidget);
    expect(find.textContaining('Клетка состоит из ядра.'), findsOneWidget);
  });

  testWidgets('notes written offline wait on the device', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 900));
    await openNotes(tester);
    app.server.offline = true;

    await newNote(tester, 'Без сети');
    await goBack(tester);
    await settle(tester);

    expect(
      find.textContaining('Правки сохранены на устройстве'),
      findsOneWidget,
    );
    expect(app.server.notes.notes, isEmpty);

    app.server.offline = false;
    await tester.pump(const Duration(seconds: 6));
    await tester.pumpAndSettle();

    expect(app.server.notes.notes, hasLength(1));
    expect(find.textContaining('Правки сохранены на устройстве'), findsNothing);
  });

  testWidgets('a [[link]] creates the target and shows a backlink', (
    tester,
  ) async {
    await pumpApp(tester, size: const Size(390, 900));
    await openNotes(tester);
    await newNote(tester, 'Подробнее в [[Митоз]]');

    await tester.tap(find.byTooltip('Связи'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Создать заметку «Митоз»'));
    await tester.pumpAndSettle();

    expect(find.text('Митоз'), findsOneWidget);
    await tester.tap(find.byTooltip('Связи'));
    await tester.pumpAndSettle();
    expect(find.text('Ссылаются сюда'), findsOneWidget);
    expect(find.text('Без названия'), findsOneWidget);
  });

  testWidgets('a note shared for reading cannot be edited', (tester) async {
    final server = FakeServer();
    server.notes.notes['n-shared'] = sharedNote('reader');
    server.notes.updates['n-shared'] = [];
    await pumpApp(tester, size: const Size(390, 900), server: server);
    await openNotes(tester);
    await settle(tester);

    expect(find.text('Общая'), findsOneWidget);
    await tester.tap(find.text('Общая'));
    await tester.pumpAndSettle();

    expect(find.text('Вы можете только читать эту заметку'), findsOneWidget);
    expect(tester.widget<TextField>(find.byType(TextField)).readOnly, isTrue);
  });

  testWidgets('the owner shares a note with a classmate', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 900));
    await openNotes(tester);
    await newNote(tester, 'Для группы');

    await tester.tap(find.byTooltip('Совместный доступ'));
    await tester.pumpAndSettle();
    expect(
      find.text('Заметкой пока никто не пользуется, кроме вас.'),
      findsOneWidget,
    );

    await tester.tap(find.text('Добавить человека'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.descendant(
        of: find.byType(AlertDialog),
        matching: find.byType(TextField),
      ),
      'Ерлан',
    );
    await settle(tester);
    await tester.tap(find.text('Ерлан Ким').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Может править'));
    await settle(tester);

    final id = app.server.notes.notes.keys.single;
    expect(app.server.notes.members[id]!.single['canEdit'], isTrue);
    expect(find.text('Ерлан Ким'), findsOneWidget);
  });

  testWidgets('a member leaves a shared note', (tester) async {
    final server = FakeServer();
    server.notes.notes['n-shared'] = sharedNote('editor');
    server.notes.updates['n-shared'] = [];
    await pumpApp(tester, size: const Size(390, 900), server: server);
    await openNotes(tester);
    await settle(tester);
    await tester.tap(find.text('Общая'));
    await tester.pumpAndSettle();
    await tester.tap(find.byTooltip('Совместный доступ'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Выйти из заметки'));
    await settle(tester);

    expect(server.notes.notes, isEmpty);
    expect(find.text('Общая'), findsNothing);
    expect(find.text('Заметки и карта поиска'), findsOneWidget);
  });

  testWidgets('notes go out as .md and come back in', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 900));
    await openNotes(tester);
    await newNote(tester, 'Текст для архива');
    await goBack(tester);

    await tester.tap(find.byType(PopupMenuButton<String>));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Скачать все заметки (.md)'));
    await tester.pumpAndSettle();

    final zip = ZipDecoder().decodeBytes(app.saver.saved.values.single);
    expect(zip.single.name, 'Без названия.md');
    expect(utf8.decode(zip.single.content), 'Текст для архива');

    app.opener.next = [
      OpenedFile(
        'Лекция.md',
        utf8.encode('---\ntags: [лекция]\n---\nО клетке'),
      ),
    ];
    await tester.tap(find.byType(PopupMenuButton<String>));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Загрузить заметки (.md или .zip)'));
    await settle(tester);

    expect(find.text('Добавлено заметок: 1'), findsOneWidget);
    expect(find.text('Лекция'), findsOneWidget);
    expect(find.textContaining('#лекция'), findsOneWidget);
  });

  testWidgets('a saved version can be restored', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 900));
    await openNotes(tester);
    await newNote(tester, 'Первая мысль');

    await tester.tap(find.byType(PopupMenuButton<String>));
    await tester.pumpAndSettle();
    await tester.tap(find.text('История версий'));
    await tester.pumpAndSettle();
    expect(find.text('Сохранённых версий пока нет'), findsOneWidget);
    await tester.tap(find.text('Сохранить версию'));
    await settle(tester);
    expect(app.server.notes.versions.values.single, hasLength(1));
    await goBack(tester);

    await tester.enterText(find.byType(TextField), 'Совсем другое');
    await settle(tester);

    await tester.tap(find.byType(PopupMenuButton<String>));
    await tester.pumpAndSettle();
    await tester.tap(find.text('История версий'));
    await tester.pumpAndSettle();
    await tester.tap(find.byIcon(Icons.history));
    await tester.pumpAndSettle();
    expect(find.text('Первая мысль'), findsOneWidget);
    await tester.tap(find.text('Вернуть эту версию'));
    await settle(tester);

    expect(find.text('Первая мысль'), findsOneWidget);
    expect(
      tester.widget<TextField>(find.byType(TextField)).controller!.text,
      'Первая мысль',
    );
  });

  testWidgets('the graph opens a note on tap', (tester) async {
    await pumpApp(tester, size: const Size(390, 900));
    await openNotes(tester);
    await newNote(tester, 'Одна');
    await goBack(tester);

    await tester.tap(find.byTooltip('Граф связей'));
    await tester.pumpAndSettle();
    expect(find.bySemanticsLabel('Граф связей, заметок: 1'), findsOneWidget);
    await tester.tapAt(tester.getCenter(find.byType(CustomPaint).last));
    await tester.pumpAndSettle();
    expect(find.text('Одна'), findsOneWidget);
  });
}
