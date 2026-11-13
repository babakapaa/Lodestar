class LeagueSeason {
  const LeagueSeason({
    required this.id,
    required this.number,
    required this.startedAt,
    required this.endedAt,
  });

  factory LeagueSeason.fromJson(Map<String, Object?> json) {
    final endedAt = json['endedAt'];
    return LeagueSeason(
      id: json['id']! as String,
      number: (json['number'] as num?)?.toInt() ?? 1,
      startedAt: DateTime.parse(json['startedAt']! as String),
      endedAt: endedAt is String ? DateTime.parse(endedAt) : null,
    );
  }

  final String id;
  final int number;
  final DateTime startedAt;
  final DateTime? endedAt;
}
