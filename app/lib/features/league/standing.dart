import 'league_display.dart';

class Standing {
  const Standing({
    required this.points,
    required this.division,
    required this.toNext,
    required this.growth,
    required this.display,
    required this.nickname,
    required this.gpuHours,
  });

  factory Standing.fromJson(Map<String, Object?> json) => Standing(
    points: (json['points'] as num?)?.toInt() ?? 0,
    division: (json['division'] as num?)?.toInt() ?? 0,
    toNext: (json['toNext'] as num?)?.toInt(),
    growth: (json['growth'] as num?)?.toInt() ?? 0,
    display: LeagueDisplay.parse(json['display']),
    nickname: json['nickname'] as String? ?? '',
    gpuHours: (json['gpuHours'] as num?)?.toInt() ?? 0,
  );

  final int points;
  final int division;
  final int? toNext;
  final int growth;
  final LeagueDisplay display;
  final String nickname;
  final int gpuHours;
}
