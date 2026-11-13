import 'navigator_intent.dart';

class Question {
  const Question({
    required this.key,
    required this.text,
    this.intent,
    this.more = false,
  });

  final String key;
  final String text;
  final NavigatorIntent? intent;
  final bool more;

  Map<String, Object?> toJson() => {
    'key': key,
    'text': text,
    'intent': ?intent?.wire,
    'more': more,
  };
}
