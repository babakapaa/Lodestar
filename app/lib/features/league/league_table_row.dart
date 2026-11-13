import 'league_display.dart';

class LeagueTableRow {
  const LeagueTableRow({
    required this.place,
    required this.name,
    required this.display,
    required this.points,
    required this.growth,
    required this.me,
    required this.topGrowth,
  });

  factory LeagueTableRow.fromJson(Map<String, Object?> json) => LeagueTableRow(
    place: (json['place'] as num?)?.toInt() ?? 0,
    name: json['name'] as String? ?? '',
    display: LeagueDisplay.parse(json['display']),
    points: (json['points'] as num?)?.toInt() ?? 0,
    growth: (json['growth'] as num?)?.toInt() ?? 0,
    me: json['me'] as bool? ?? false,
    topGrowth: json['topGrowth'] as bool? ?? false,
  );

  final int place;
  final String name;
  final LeagueDisplay display;
  final int points;
  final int growth;
  final bool me;
  final bool topGrowth;
}
