import 'dart:convert';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/app/lodestar_app.dart';
import 'package:lodestar/features/notes/canvas/canvas_card_view.dart';
import 'package:lodestar/features/notes/canvas/canvas_view.dart';
import 'package:lodestar/features/notes/notes_providers.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';
import 'notes_test.dart' show openNotes, settle;

void main() {
  testWidgets('a canvas gets cards that are saved as JSON Canvas', (
    tester,
  ) async {
    await pumpApp(tester, size: const Size(1440, 900));
    await openNotes(tester);
    await tester.tap(find.byTooltip('Новый холст'));
    await settle(tester);
    expect(find.byType(CanvasView), findsOneWidget);
    expect(find.text('Холст'), findsWidgets);

    await tester.tap(find.byTooltip('Добавить карточку'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.descendant(
        of: find.byType(CanvasCardView),
        matching: find.byType(TextField),
      ),
      'Первая мысль',
    );
    await settle(tester);

    final container = ProviderScope.containerOf(
      tester.element(find.byType(LodestarApp)),
    );
    final note = container
        .read(notesProvider)
        .value!
        .singleWhere((note) => note.meta.isCanvas);
    expect(note.meta.path, 'Холст.canvas');
    final json = jsonDecode(note.body) as Map<String, Object?>;
    final nodes = json['nodes']! as List<Object?>;
    expect((nodes.single! as Map<String, Object?>)['text'], 'Первая мысль');
  });
}
