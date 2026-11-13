class TeamScope {
  const TeamScope({
    required this.kind,
    required this.id,
    required this.title,
    required this.context,
    required this.link,
    required this.minSize,
    required this.maxSize,
    required this.selfFormed,
    required this.freeLeaveUntil,
    required this.dueAt,
    required this.archived,
  });

  factory TeamScope.fromJson(Map<String, Object?> json) {
    final freeLeaveUntil = json['freeLeaveUntil'];
    final dueAt = json['dueAt'];
    return TeamScope(
      kind: json['kind']! as String,
      id: json['id']! as String,
      title: json['title'] as String? ?? '',
      context: json['context'] as String? ?? '',
      link: json['link'] as String? ?? '',
      minSize: json['minSize'] as int? ?? 2,
      maxSize: json['maxSize'] as int? ?? 2,
      selfFormed: json['selfFormed'] as bool? ?? false,
      freeLeaveUntil: freeLeaveUntil is String
          ? DateTime.parse(freeLeaveUntil)
          : null,
      dueAt: dueAt is String ? DateTime.parse(dueAt) : null,
      archived: json['archived'] as bool? ?? false,
    );
  }

  final String kind;
  final String id;
  final String title;
  final String context;
  final String link;
  final int minSize;
  final int maxSize;
  final bool selfFormed;
  final DateTime? freeLeaveUntil;
  final DateTime? dueAt;
  final bool archived;
}
