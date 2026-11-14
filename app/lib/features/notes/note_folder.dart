import 'note.dart';

class NoteFolder {
  NoteFolder(this.name, this.path);

  factory NoteFolder.of(Iterable<Note> notes) {
    final root = NoteFolder('', '');
    for (final note in notes) {
      if (note.deleted) {
        continue;
      }
      var folder = root;
      final parts = note.meta.folder.isEmpty
          ? const <String>[]
          : note.meta.folder.split('/');
      for (final part in parts) {
        folder = folder.child(part);
      }
      folder.notes.add(note);
    }
    root._sort();
    return root;
  }

  final String name;
  final String path;
  final List<NoteFolder> folders = [];
  final List<Note> notes = [];

  NoteFolder child(String name) {
    for (final folder in folders) {
      if (folder.name == name) {
        return folder;
      }
    }
    final folder = NoteFolder(name, path.isEmpty ? name : '$path/$name');
    folders.add(folder);
    return folder;
  }

  void _sort() {
    folders.sort(
      (a, b) => a.name.toLowerCase().compareTo(b.name.toLowerCase()),
    );
    notes.sort(
      (a, b) => a.meta.name.toLowerCase().compareTo(b.meta.name.toLowerCase()),
    );
    for (final folder in folders) {
      folder._sort();
    }
  }
}
