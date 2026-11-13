import 'setup_legal_answer.dart';
import 'setup_problem.dart';
import 'setup_section.dart';

class SetupOverview {
  const SetupOverview({
    required this.launched,
    required this.protocolVersion,
    required this.signed,
    required this.sections,
    required this.legal,
    required this.problems,
    this.pilot = const [],
  });

  factory SetupOverview.fromJson(Map<String, Object?> json) => SetupOverview(
    launched: json['launched'] as bool? ?? false,
    protocolVersion: json['protocolVersion'] as int? ?? 0,
    signed: json['signedAt'] != null,
    sections: [
      for (final item in json['sections'] as List<Object?>? ?? const [])
        SetupSection.fromJson(item! as Map<String, Object?>),
    ],
    legal: [
      for (final item in json['legal'] as List<Object?>? ?? const [])
        SetupLegalAnswer.fromJson(item! as Map<String, Object?>),
    ],
    problems: switch (json['problems']) {
      final List<Object?> items => [
        for (final item in items)
          SetupProblem.fromJson(item! as Map<String, Object?>),
      ],
      _ => null,
    },
    pilot: [
      for (final item in json['pilot'] as List<Object?>? ?? const [])
        item! as String,
    ],
  );

  final bool launched;
  final int protocolVersion;
  final bool signed;
  final List<SetupSection> sections;
  final List<SetupLegalAnswer> legal;
  final List<SetupProblem>? problems;
  final List<String> pilot;

  SetupSection? section(String key) {
    for (final s in sections) {
      if (s.key == key) {
        return s;
      }
    }
    return null;
  }
}
