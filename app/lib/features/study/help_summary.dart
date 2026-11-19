import 'help_steps.dart';

class HelpSummary {
  const HelpSummary({
    required this.replies,
    required this.steps,
    required this.sources,
    required this.articles,
  });

  factory HelpSummary.fromJson(Map<String, Object?> json) => HelpSummary(
    replies: json['replies'] as int? ?? 0,
    steps: HelpSteps.fromJson(
      json['steps'] as Map<String, Object?>? ?? const {},
    ),
    sources: _strings(json['sources']),
    articles: _strings(json['articles']),
  );

  static HelpSummary? maybe(Object? json) =>
      json is Map<String, Object?> ? HelpSummary.fromJson(json) : null;

  static List<String> _strings(Object? json) => [
    for (final item in json as List<Object?>? ?? const []) item! as String,
  ];

  final int replies;
  final HelpSteps steps;
  final List<String> sources;
  final List<String> articles;
}
