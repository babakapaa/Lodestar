class AuditEntry {
  const AuditEntry({
    required this.id,
    required this.actorId,
    required this.action,
    required this.targetKind,
    required this.targetId,
    required this.details,
    required this.at,
  });

  factory AuditEntry.fromJson(Map<String, Object?> json) {
    final details = json['details'];
    return AuditEntry(
      id: (json['id']! as num).toInt(),
      actorId: json['actorId'] as String?,
      action: json['action']! as String,
      targetKind: json['targetKind']! as String,
      targetId: json['targetId']! as String,
      details: details is Map<String, Object?> ? details : const {},
      at: DateTime.parse(json['at']! as String),
    );
  }

  final int id;
  final String? actorId;
  final String action;
  final String targetKind;
  final String targetId;
  final Map<String, Object?> details;
  final DateTime at;
}
