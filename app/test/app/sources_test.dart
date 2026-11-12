import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import 'app_harness.dart';
import 'notes_test.dart' show openNotes;

Future<void> openSources(WidgetTester tester) async {
  await openNotes(tester);
  await tester.tap(find.byType(PopupMenuButton<String>));
  await tester.pumpAndSettle();
  await tester.tap(find.text('Источники'));
  await tester.pumpAndSettle();
}

Future<void> save(WidgetTester tester) async {
  await tester.dragUntilVisible(
    find.text('Сохранить'),
    find.byType(ListView).last,
    const Offset(0, -300),
  );
  await tester.tap(find.text('Сохранить'));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('a source found by DOI is saved as a GOST reference', (
    tester,
  ) async {
    final app = await pumpApp(tester, size: const Size(390, 900));
    await openSources(tester);
    expect(find.text('Источников пока нет'), findsOneWidget);

    await tester.tap(find.text('Найти источник'));
    await tester.pumpAndSettle();
    await tester.enterText(find.byType(TextField), 'doi.org/10.1/none');
    await tester.tap(find.text('Найти'));
    await tester.pumpAndSettle();
    expect(
      find.text('Не удалось найти описание. Добавьте источник вручную.'),
      findsOneWidget,
    );

    await tester.enterText(find.byType(TextField), '10.1000/xyz');
    await tester.tap(find.text('Найти'));
    await tester.pumpAndSettle();
    expect(find.text('Новый источник'), findsOneWidget);
    expect(find.textContaining('Иванов, А. Б. Строение клетки'), findsWidgets);

    await save(tester);
    final saved = app.server.sources.sources.values.single;
    expect(saved['title'], 'Строение клетки');
    expect(find.textContaining('Вестник биологии'), findsOneWidget);
    expect(
      find.text('Оформлено автоматически, библиотекой не проверено'),
      findsOneWidget,
    );

    await tester.tap(find.textContaining('Вестник биологии'));
    await tester.pumpAndSettle();
    await tester.tap(find.byTooltip('Удалить источник'));
    await tester.pumpAndSettle();
    expect(app.server.sources.sources, isEmpty);
    expect(find.text('Источников пока нет'), findsOneWidget);
  });

  testWidgets('a source is added by hand', (tester) async {
    final app = await pumpApp(tester, size: const Size(390, 900));
    await openSources(tester);
    await tester.tap(find.byTooltip('Добавить вручную'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.widgetWithText(TextField, 'Название'),
      'Биология клетки',
    );
    await tester.enterText(
      find.widgetWithText(
        TextField,
        'Авторы, по одному в строке: Фамилия, Имя Отчество',
      ),
      'Ким, Е.',
    );
    await tester.pumpAndSettle();
    await save(tester);
    expect(app.server.sources.sources.values.single['authors'], ['Ким, Е.']);
    expect(find.textContaining('Ким, Е. Биология клетки'), findsOneWidget);
  });
}
