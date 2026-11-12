import 'note_meta.dart';
import 'note_role.dart';

class RemoteNote {
  const RemoteNote({
    required this.id,
    required this.ownerId,
    required this.role,
    required this.meta,
    required this.seq,
    required this.changedAt,
    required this.deletedAt,
  });

  factory RemoteNote.fromJson(Map<String, Object?> json) => RemoteNote(
    id: json['id']! as String,
    ownerId: json['ownerId']! as String,
    role: NoteRole.parse(json['role']! as String),
    meta: NoteMeta.fromJson(json['meta']! as Map<String, Object?>),
    seq: (json['seq']! as num).toInt(),
    changedAt: DateTime.parse(json['changedAt']! as String),
    deletedAt: switch (json['deletedAt']) {
      final String at => DateTime.parse(at),
      _ => null,
    },
  );

  final String id;
  final String ownerId;
  final NoteRole role;
  final NoteMeta meta;
  final int seq;
  final DateTime changedAt;
  final DateTime? deletedAt;
}
