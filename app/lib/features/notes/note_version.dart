class NoteVersion {
  const NoteVersion({
    required this.id,
    required this.seq,
    required this.bytes,
    required this.createdAt,
  });

  factory NoteVersion.fromJson(Map<String, Object?> json) => NoteVersion(
    id: json['id']! as String,
    seq: (json['seq']! as num).toInt(),
    bytes: (json['bytes']! as num).toInt(),
    createdAt: DateTime.parse(json['createdAt']! as String),
  );

  final String id;
  final int seq;
  final int bytes;
  final DateTime createdAt;
}
