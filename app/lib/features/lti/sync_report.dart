class SyncReport {
  const SyncReport({
    required this.members,
    required this.enrolled,
    required this.works,
    required this.scores,
  });

  factory SyncReport.fromJson(Map<String, Object?> json) => SyncReport(
    members: (json['members'] as num?)?.toInt() ?? 0,
    enrolled: (json['enrolled'] as num?)?.toInt() ?? 0,
    works: (json['works'] as num?)?.toInt() ?? 0,
    scores: (json['scores'] as num?)?.toInt() ?? 0,
  );

  final int members;
  final int enrolled;
  final int works;
  final int scores;
}
