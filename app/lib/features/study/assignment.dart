import 'grading.dart';
import 'map_rules.dart';
import 'navigator_mode.dart';

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
}
