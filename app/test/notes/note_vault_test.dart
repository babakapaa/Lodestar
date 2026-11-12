import 'dart:convert';
import 'dart:typed_data';

import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/files/opened_file.dart';
import 'package:lodestar/features/notes/front_matter.dart';
import 'package:lodestar/features/notes/note.dart';
import 'package:lodestar/features/notes/note_meta.dart';
import 'package:lodestar/features/notes/note_role.dart';
import 'package:lodestar/features/notes/note_vault.dart';
import 'package:lodestar/features/notes/notes_graph.dart';

Note note(String path, String body, {List<String> tags = const []}) => Note(
  id: path,
  role: NoteRole.owner,
  meta: NoteMeta(path: path, tags: tags, links: const []),
  body: body,
  changedAt: DateTime.utc(2026, 11, 12),
  deletedAt: null,
  error: null,
);

void main() {
  test('reads Obsidian tags and keeps other properties', () {
    final parsed = FrontMatter.parse(
      '---\naliases: [Клетка]\ntags:\n  - биология\n  - "#экзамен"\n---\nТекст',
    );
    expect(parsed.tags, ['биология', 'экзамен']);
    expect(parsed.body, '---\naliases: [Клетка]\n---\nТекст');

    final inline = FrontMatter.parse('---\ntags: [a, b]\n---\nТекст');
    expect(inline.tags, ['a', 'b']);
    expect(inline.body, 'Текст');
  });

  test('a vault survives export and import', () {
    final zip = NoteVault.export([
      note('Биология/Клетка', 'См. [[Митоз]]', tags: ['экзамен']),
      note('Митоз', 'Деление'),
      note('Вопрос: что?', 'x'),
    ]);

    final notes = NoteVault.read([OpenedFile('vault.zip', zip)]);

    expect(
      [for (final n in notes) n.path],
      ['Биология/Клетка', 'Митоз', 'Вопрос- что-'],
    );
    expect(notes.first.tags, ['экзамен']);
    expect(notes.first.body, 'См. [[Митоз]]');
  });

  test('single markdown files are imported, hidden folders are not', () {
    final notes = NoteVault.read([
      OpenedFile('Лекция.md', utf8.encode('Текст')),
      OpenedFile('.obsidian/app.md', utf8.encode('x')),
      OpenedFile('picture.png', Uint8List.fromList([1, 2])),
    ]);
    expect(notes.single.path, 'Лекция');
  });

  test('the graph links notes by name and spreads them apart', () {
    final graph = NotesGraph.of([
      note('Биология/Клетка', ''),
      Note(
        id: 'm',
        role: NoteRole.owner,
        meta: const NoteMeta(path: 'Митоз', links: ['Клетка', 'Нет такой']),
        body: '',
        changedAt: DateTime.utc(2026),
        deletedAt: null,
        error: null,
      ),
      note('Одна', ''),
    ]);
    expect(graph.edges, [(0, 1)]);
    expect(graph.nodes[0].degree, 1);
    final a = graph.nodes[0].position;
    final b = graph.nodes[2].position;
    expect((a - b).distance, greaterThan(40));
    expect(a.dx.isFinite && a.dy.isFinite, isTrue);
  });
}
