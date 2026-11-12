import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/notes/note_links.dart';
import 'package:lodestar/features/notes/text_edit.dart';

void main() {
  test('finds the changed span between two texts', () {
    final edit = TextEdit.between('Клетка ядро', 'Клетка и ядро')!;
    expect((edit.index, edit.remove, edit.insert), (7, 0, 'и '));
    expect(TextEdit.between('same', 'same'), isNull);
  });

  test('never splits an emoji in half', () {
    final edit = TextEdit.between('a🧬b', 'a🧫b')!;
    expect((edit.index, edit.remove, edit.insert), (1, 2, '🧫'));
  });

  test('shifts a cursor around a remote change', () {
    const edit = TextEdit(2, 3, 'xy');
    expect(edit.shift(1), 1);
    expect(edit.shift(4), 4);
    expect(edit.shift(7), 6);
  });

  test('reads Obsidian links with aliases and headings', () {
    expect(
      NoteLinks.parse(
        '[[Клетка]], [[Биология/Ткани|ткани]], [[Митоз#Фазы]], '
        '[[Клетка]]',
      ),
      ['Клетка', 'Биология/Ткани', 'Митоз'],
    );
  });
}
