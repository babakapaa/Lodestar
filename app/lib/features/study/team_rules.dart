class TeamRules {
  const TeamRules({
    required this.minSize,
    required this.maxSize,
    required this.selfFormed,
    this.freeLeaveUntil,
  });

  factory TeamRules.fromJson(Map<String, Object?> json) {
    final until = json['freeLeaveUntil'];
    return TeamRules(
      minSize: (json['minSize'] as num?)?.toInt() ?? 2,
      maxSize: (json['maxSize'] as num?)?.toInt() ?? 2,
      selfFormed: json['selfFormed'] as bool? ?? false,
      freeLeaveUntil: until is String ? DateTime.parse(until) : null,
    );
  }

  static TeamRules? maybe(Object? json) =>
      json is Map<String, Object?> ? TeamRules.fromJson(json) : null;

  static const maxTeamSize = 12;

  final int minSize;
  final int maxSize;
  final bool selfFormed;
  final DateTime? freeLeaveUntil;

  bool get valid =>
      minSize >= 2 && maxSize >= minSize && maxSize <= maxTeamSize;

  Map<String, Object?> toJson() => {
    'minSize': minSize,
    'maxSize': maxSize,
    'selfFormed': selfFormed,
    'freeLeaveUntil': freeLeaveUntil?.toUtc().toIso8601String(),
  };
}
