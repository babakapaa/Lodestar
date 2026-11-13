class AwardGrant {
  const AwardGrant({
    required this.id,
    required this.title,
    required this.gpuHours,
    required this.userId,
    required this.reason,
    required this.grantedAt,
  });

  factory AwardGrant.fromJson(Map<String, Object?> json) => AwardGrant(
    id: json['id']! as String,
    title: json['title'] as String? ?? '',
    gpuHours: (json['gpuHours'] as num?)?.toInt() ?? 0,
    userId: json['userId'] as String? ?? '',
    reason: json['reason'] as String? ?? '',
    grantedAt: DateTime.parse(json['grantedAt']! as String),
  );

  final String id;
  final String title;
  final int gpuHours;
  final String userId;
  final String reason;
  final DateTime grantedAt;
}
