import 'dart:convert';
import 'dart:typed_data';

import 'package:archive/archive.dart';

import '../../core/files/opened_file.dart';
import 'front_matter.dart';
import 'note.dart';
import 'note_meta.dart';
import 'vault_note.dart';

abstract final class NoteVault {
  static final _unsafe = RegExp(r'[\\:*?"<>|]');

  static Uint8List export(List<Note> notes) {
    final archive = Archive();
    final used = <String>{};
    for (final note in notes) {
      if (note.deleted) {
        continue;
      }
      final base = note.meta.path.replaceAll(_unsafe, '-');
      var name = '$base.md';
      for (var n = 2; used.contains(name.toLowerCase()); n++) {
        name = '$base $n.md';
      }
      used.add(name.toLowerCase());
      archive.add(
        ArchiveFile.bytes(
          name,
          utf8.encode(FrontMatter.render(note.meta.tags, note.body)),
        ),
      );
    }
    return Uint8List.fromList(ZipEncoder().encode(archive));
  }

  static List<VaultNote> read(List<OpenedFile> files) {
    final notes = <VaultNote>[];
    for (final file in files) {
      if (file.name.toLowerCase().endsWith('.zip')) {
        for (final entry in ZipDecoder().decodeBytes(file.bytes)) {
          if (entry.isFile && _isNote(entry.name)) {
            notes.add(_note(entry.name, entry.content));
          }
        }
      } else if (_isNote(file.name)) {
        notes.add(_note(file.name, file.bytes));
      }
    }
    return notes;
  }

  static bool _isNote(String name) {
    final parts = name.split('/');
    return name.toLowerCase().endsWith('.md') &&
        !parts.any((part) => part.startsWith('.') || part == '__MACOSX');
  }

  static VaultNote _note(String name, List<int> bytes) {
    final parsed = FrontMatter.parse(
      utf8.decode(bytes, allowMalformed: true).replaceAll('\r\n', '\n'),
    );
    return VaultNote(
      path: NoteMeta.normalizePath(name.substring(0, name.length - 3)),
      tags: parsed.tags,
      body: parsed.body,
    );
  }
}
