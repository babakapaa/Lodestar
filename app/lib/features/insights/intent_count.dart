import '../navigator/navigator_intent.dart';

class IntentCount {
  const IntentCount({required this.intent, required this.students});

  factory IntentCount.fromJson(Map<String, Object?> json) => IntentCount(
    intent: NavigatorIntent.parse(json['intent']),
    students: json['students'] as int? ?? 0,
  );

  final NavigatorIntent? intent;
  final int students;
}
