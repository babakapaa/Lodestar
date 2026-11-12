import 'package:flutter_test/flutter_test.dart';
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
    await tester.tap(find.text('Название и папка'));
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
}
