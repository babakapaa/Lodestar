import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/notes/note_editor.dart';
import 'package:lodestar/features/notes/note_inspector.dart';
import 'package:lodestar/features/notes/notes_sidebar.dart';
import 'package:lodestar/features/notes/notes_tab_chip.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';
import 'notes_test.dart' show openNotes, settle;

Future<void> typeInEditor(WidgetTester tester, String text) async {
  await tester.enterText(
    find.descendant(
      of: find.byType(NoteEditor),
      matching: find.byType(TextField),
    ),
    text,
  );
  await settle(tester);
}

Future<void> press(WidgetTester tester, LogicalKeyboardKey key) async {
  await tester.sendKeyDownEvent(LogicalKeyboardKey.controlLeft);
  await tester.sendKeyEvent(key);
  await tester.sendKeyUpEvent(LogicalKeyboardKey.controlLeft);
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('wide notes open in tabs with tree, outline and links', (
    tester,
  ) async {
    await pumpApp(tester, size: const Size(1440, 900));
    await openNotes(tester);
    expect(find.byType(NotesSidebar), findsOneWidget);

    await tester.tap(find.widgetWithText(FilledButton, 'Новая заметка'));
    await settle(tester);
    expect(find.byType(NotesTabChip), findsOneWidget);
    await typeInEditor(tester, '# Клетка\n\nСм. [[Ядро]]');
    expect(
      find.descendant(
        of: find.byType(NoteInspector),
        matching: find.text('Клетка'),
      ),
      findsOneWidget,
    );
    expect(
      find.descendant(
        of: find.byType(NoteInspector),
        matching: find.text('Ядро'),
      ),
      findsOneWidget,
    );

    await press(tester, LogicalKeyboardKey.keyN);
    await settle(tester);
    expect(find.byType(NotesTabChip), findsNWidgets(2));
    await typeInEditor(tester, 'Связь с [[Без');
    expect(
      find.text('Enter вставит ссылку, Esc скроет список'),
      findsOneWidget,
    );
    await tester.sendKeyEvent(LogicalKeyboardKey.enter);
    await settle(tester);
    final editor = tester.widget<TextField>(
      find.descendant(
        of: find.byType(NoteEditor),
        matching: find.byType(TextField),
      ),
    );
    expect(editor.controller!.text, 'Связь с [[Без названия]]');

    await press(tester, LogicalKeyboardKey.keyW);
    expect(find.byType(NotesTabChip), findsOneWidget);
  });

  testWidgets('quick switcher creates a missing note', (tester) async {
    await pumpApp(tester, size: const Size(1440, 900));
    await openNotes(tester);
    await press(tester, LogicalKeyboardKey.keyO);
    await tester.enterText(find.byType(TextField).last, 'Митоз');
    await tester.pumpAndSettle();
    await tester.tap(find.text('Создать заметку «Митоз»'));
    await settle(tester);
    expect(
      find.descendant(
        of: find.byType(NotesTabChip),
        matching: find.text('Митоз'),
      ),
      findsOneWidget,
    );
  });
}
