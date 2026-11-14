import 'note.dart';

abstract final class NoteResolver {
  static Note? resolve(String link, Iterable<Note> notes) {
    for (final candidate in notes) {
      if (!candidate.deleted &&
          (candidate.meta.path == link || candidate.meta.name == link)) {
        return candidate;
      }
    }
    return null;
  }

  static List<Note> backlinks(Note note, Iterable<Note> notes) => [
    for (final other in notes)
      if (other.id != note.id && !other.deleted && other.linksTo(note)) other,
  ];

  static List<Note> search(String query, Iterable<Note> notes) {
    final needle = query.trim().toLowerCase();
    final scored = <(int, Note)>[];
    for (final note in notes) {
      if (note.deleted) {
        continue;
      }
      final name = note.meta.name.toLowerCase();
      final path = note.meta.path.toLowerCase();
      final score = needle.isEmpty
          ? 3
          : name.startsWith(needle)
          ? 0
          : name.contains(needle)
          ? 1
          : path.contains(needle)
          ? 2
          : -1;
      if (score >= 0) {
        scored.add((score, note));
      }
    }
    scored.sort((a, b) {
      final byScore = a.$1.compareTo(b.$1);
      return byScore != 0
          ? byScore
          : needle.isEmpty
          ? b.$2.changedAt.compareTo(a.$2.changedAt)
          : a.$2.meta.name.compareTo(b.$2.meta.name);
    });
    return [for (final (_, note) in scored) note];
  }
}
