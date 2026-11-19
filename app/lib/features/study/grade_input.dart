import 'grade_parts.dart';

class GradeInput {
  const GradeInput({
    required this.points,
    required this.passed,
    required this.comment,
    this.parts,
  });

  final double? points;
  final bool? passed;
  final GradeParts? parts;
  final String comment;

  Map<String, Object?> toJson() => {
    'points': points,
    'passed': passed,
    'parts': ?parts?.toJson(),
    'comment': comment,
  };
}
