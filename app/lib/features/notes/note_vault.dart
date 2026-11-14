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
      final canvas = note.meta.isCanvas;
      final ext = canvas ? NoteMeta.canvasExtension : '.md';
      final path = note.meta.path.replaceAll(_unsafe, '-');
      final base = canvas ? path.substring(0, path.length - ext.length) : path;
      var name = '$base$ext';
      for (var n = 2; used.contains(name.toLowerCase()); n++) {
        name = '$base $n$ext';
      }
      used.add(name.toLowerCase());
      archive.add(
        ArchiveFile.bytes(
          name,
          utf8.encode(
            canvas ? note.body : FrontMatter.render(note.meta.tags, note.body),
          ),
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
    final lower = name.toLowerCase();
    return (lower.endsWith('.md') ||
            lower.endsWith(NoteMeta.canvasExtension)) &&
        !parts.any((part) => part.startsWith('.') || part == '__MACOSX');
  }

  static VaultNote _note(String name, List<int> bytes) {
    final text = utf8
        .decode(bytes, allowMalformed: true)
        .replaceAll('\r\n', '\n');
    if (name.toLowerCase().endsWith(NoteMeta.canvasExtension)) {
      return VaultNote(
        path: NoteMeta.normalizePath(name),
        tags: const [],
        body: text,
      );
    }
    final parsed = FrontMatter.parse(text);
    return VaultNote(
      path: NoteMeta.normalizePath(name.substring(0, name.length - 3)),
      tags: parsed.tags,
      body: parsed.body,
    );
  }
}
