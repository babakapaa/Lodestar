import 'defense_format.dart';
import 'grading.dart';
import 'map_rules.dart';
import 'navigator_mode.dart';
import 'team_rules.dart';

class Assignment {
  const Assignment({
    required this.id,
    required this.courseId,
    required this.title,
    required this.description,
    required this.dueAt,
    required this.grading,
    required this.maxPoints,
    required this.published,
    this.searchMap,
    this.navigatorMode = NavigatorMode.course,
    this.navigatorSummary = false,
    this.defense = DefenseFormat.none,
    this.team,
  });

  factory Assignment.fromJson(Map<String, Object?> json) {
    final dueAt = json['dueAt'];
    return Assignment(
      id: json['id']! as String,
      courseId: json['courseId']! as String,
      title: json['title']! as String,
      description: json['description'] as String? ?? '',
      dueAt: dueAt is String ? DateTime.parse(dueAt) : null,
      grading: Grading.parse(json['grading']),
      maxPoints: (json['maxPoints'] as num?)?.toDouble(),
      published: json['publishedAt'] != null,
      searchMap: MapRules.maybe(json['searchMap']),
      navigatorMode: NavigatorMode.parse(json['navigatorMode']),
      navigatorSummary: json['navigatorSummary'] as bool? ?? false,
      defense: DefenseFormat.parse(json['defense']),
      team: TeamRules.maybe(json['team']),
    );
  }

  final String id;
  final String courseId;
  final String title;
  final String description;
  final DateTime? dueAt;
  final Grading grading;
  final double? maxPoints;
  final bool published;
  final MapRules? searchMap;
  final NavigatorMode navigatorMode;
  final bool navigatorSummary;
  final DefenseFormat defense;
  final TeamRules? team;
}
