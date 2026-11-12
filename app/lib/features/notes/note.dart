import 'notes_database.dart';
import 'note_meta.dart';
import 'note_role.dart';

class Note {
  const Note({
    required this.id,
    required this.role,
    required this.meta,
    required this.body,
    required this.changedAt,
    required this.deletedAt,
    required this.error,
  });

  factory Note.fromRow(NoteRow row) => Note(
    id: row.id,
    role: NoteRole.parse(row.role),
    meta: NoteMeta(
      path: row.path,
      tags: row.tags,
      links: row.links,
      assignmentId: row.assignmentId,
      lessonId: row.lessonId,
    ),
    body: row.body,
    changedAt: row.changedAt,
    deletedAt: row.deletedAt,
    error: row.error,
  );

  final String id;
  final NoteRole role;
  final NoteMeta meta;
  final String body;
  final DateTime changedAt;
  final DateTime? deletedAt;
  final String? error;

  bool get deleted => deletedAt != null;

  bool linksTo(Note other) => meta.links.any(
    (link) => link == other.meta.name || link == other.meta.path,
  );
}
