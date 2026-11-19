import 'defense_format.dart';
import 'grading.dart';
import 'map_rules.dart';
import 'navigator_mode.dart';
import 'team_rules.dart';

class AssignmentInput {
  const AssignmentInput({
    required this.title,
    required this.description,
    required this.dueAt,
    required this.grading,
    required this.maxPoints,
    required this.published,
    required this.files,
    this.searchMap,
    this.navigatorMode = NavigatorMode.course,
    this.navigatorSummary = false,
    this.defense = DefenseFormat.none,
    this.team,
  });

  final String title;
  final String description;
  final DateTime? dueAt;
  final Grading grading;
  final double? maxPoints;
  final bool published;
  final List<String> files;
  final MapRules? searchMap;
  final NavigatorMode navigatorMode;
  final bool navigatorSummary;
  final DefenseFormat defense;
  final TeamRules? team;

  Map<String, Object?> toJson() => {
    'title': title,
    'description': description,
    'dueAt': dueAt?.toUtc().toIso8601String(),
    'grading': grading.name,
    'maxPoints': grading == Grading.points ? maxPoints : null,
    'published': published,
    'files': files,
    'searchMap': searchMap?.toJson(),
    'navigatorMode': navigatorMode.wire,
    'navigatorSummary': navigatorSummary,
    'defense': defense.wire,
    'team': team?.toJson(),
  };
}
