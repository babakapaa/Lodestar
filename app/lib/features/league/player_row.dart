import '../study/course_person.dart';
import 'league_display.dart';

class PlayerRow {
  const PlayerRow({
    required this.person,
    required this.division,
    required this.points,
    required this.growth,
    required this.topGrowth,
    required this.display,
    required this.nickname,
  });

  factory PlayerRow.fromJson(Map<String, Object?> json) => PlayerRow(
    person: CoursePerson.fromJson(json['person']! as Map<String, Object?>),
    division: (json['division'] as num?)?.toInt() ?? 0,
    points: (json['points'] as num?)?.toInt() ?? 0,
    growth: (json['growth'] as num?)?.toInt() ?? 0,
    topGrowth: json['topGrowth'] as bool? ?? false,
    display: LeagueDisplay.parse(json['display']),
    nickname: json['nickname'] as String? ?? '',
  );

  final CoursePerson person;
  final int division;
  final int points;
  final int growth;
  final bool topGrowth;
  final LeagueDisplay display;
  final String nickname;
}
