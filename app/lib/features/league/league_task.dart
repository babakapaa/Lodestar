class LeagueTask {
  const LeagueTask({
    required this.id,
    required this.leagueId,
    required this.title,
    required this.about,
    required this.partner,
    required this.points,
    required this.minSize,
    required this.maxSize,
    required this.selfFormed,
    required this.dueAt,
    required this.closed,
  });

  factory LeagueTask.fromJson(Map<String, Object?> json) {
    final dueAt = json['dueAt'];
    return LeagueTask(
      id: json['id']! as String,
      leagueId: json['leagueId']! as String,
      title: json['title'] as String? ?? '',
      about: json['about'] as String? ?? '',
      partner: json['partner'] as String? ?? '',
      points: (json['points'] as num?)?.toInt() ?? 0,
      minSize: (json['minSize'] as num?)?.toInt() ?? 2,
      maxSize: (json['maxSize'] as num?)?.toInt() ?? 2,
      selfFormed: json['selfFormed'] as bool? ?? true,
      dueAt: dueAt is String ? DateTime.parse(dueAt) : null,
      closed: json['closedAt'] != null,
    );
  }

  final String id;
  final String leagueId;
  final String title;
  final String about;
  final String partner;
  final int points;
  final int minSize;
  final int maxSize;
  final bool selfFormed;
  final DateTime? dueAt;
  final bool closed;
}
