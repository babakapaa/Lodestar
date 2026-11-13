import 'entry_kind.dart';

class LeagueEntry {
  const LeagueEntry({
    required this.id,
    required this.kind,
    required this.title,
    required this.score,
    required this.max,
    required this.weight,
    required this.points,
    required this.reason,
    required this.createdAt,
  });

  factory LeagueEntry.fromJson(Map<String, Object?> json) => LeagueEntry(
    id: json['id']! as String,
    kind: EntryKind.parse(json['kind']),
    title: json['title'] as String? ?? '',
    score: (json['score'] as num?)?.toDouble(),
    max: (json['max'] as num?)?.toDouble(),
    weight: (json['weight'] as num?)?.toInt(),
    points: (json['points'] as num?)?.toInt() ?? 0,
    reason: json['reason'] as String? ?? '',
    createdAt: DateTime.parse(json['createdAt']! as String),
  );

  final String id;
  final EntryKind kind;
  final String title;
  final double? score;
  final double? max;
  final int? weight;
  final int points;
  final String reason;
  final DateTime createdAt;
}
