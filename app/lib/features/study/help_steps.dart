class HelpSteps {
  const HelpSteps({
    required this.question,
    required this.direction,
    required this.source,
    required this.fragment,
    required this.example,
  });

  factory HelpSteps.fromJson(Map<String, Object?> json) => HelpSteps(
    question: json['question'] as int? ?? 0,
    direction: json['direction'] as int? ?? 0,
    source: json['source'] as int? ?? 0,
    fragment: json['fragment'] as int? ?? 0,
    example: json['example'] as int? ?? 0,
  );

  final int question;
  final int direction;
  final int source;
  final int fragment;
  final int example;
}
