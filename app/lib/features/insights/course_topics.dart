import 'intent_count.dart';
import 'topic_count.dart';
import 'week_count.dart';

class CourseTopics {
  const CourseTopics({
    required this.threshold,
    required this.students,
    required this.topics,
    required this.intents,
    required this.trend,
    required this.suggestion,
  });

  factory CourseTopics.fromJson(Map<String, Object?> json) {
    List<T> list<T>(String key, T Function(Map<String, Object?>) parse) => [
      for (final item in json[key] as List<Object?>? ?? const [])
        parse(item! as Map<String, Object?>),
    ];
    final suggestion = json['suggestion'] as Map<String, Object?>?;
    return CourseTopics(
      threshold: json['threshold'] as int? ?? 5,
      students: json['students'] as int?,
      topics: list('topics', TopicCount.fromJson),
      intents: list('intents', IntentCount.fromJson),
      trend: list('trend', WeekCount.fromJson),
      suggestion: suggestion == null ? null : TopicCount.fromJson(suggestion),
    );
  }

  final int threshold;
  final int? students;
  final List<TopicCount> topics;
  final List<IntentCount> intents;
  final List<WeekCount> trend;
  final TopicCount? suggestion;
}
