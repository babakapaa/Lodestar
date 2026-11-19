import 'article.dart';
import 'citation.dart';
import 'hint_step.dart';
import 'navigator_intent.dart';

class NavigatorMessage {
  const NavigatorMessage({
    required this.seq,
    required this.fromNavigator,
    required this.text,
    required this.intent,
    required this.step,
    required this.lastStep,
    required this.support,
    required this.citations,
    this.articles = const [],
  });

  factory NavigatorMessage.fromJson(Map<String, Object?> json) =>
      NavigatorMessage(
        seq: json['seq']! as int,
        fromNavigator: json['role'] == 'navigator',
        text: json['text'] as String? ?? '',
        intent: NavigatorIntent.parse(json['intent']),
        step: HintStep.parse(json['step']),
        lastStep: json['lastStep'] as bool? ?? false,
        support: json['support'] as String? ?? '',
        citations: [
          for (final item in json['citations'] as List<Object?>? ?? const [])
            Citation.fromJson(item! as Map<String, Object?>),
        ],
        articles: [
          for (final item in json['articles'] as List<Object?>? ?? const [])
            Article.fromJson(item! as Map<String, Object?>),
        ],
      );

  final int seq;
  final bool fromNavigator;
  final String text;
  final NavigatorIntent? intent;
  final HintStep? step;
  final bool lastStep;
  final String support;
  final List<Citation> citations;
  final List<Article> articles;
}
