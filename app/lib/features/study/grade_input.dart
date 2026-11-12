class GradeInput {
  const GradeInput({
    required this.points,
    required this.passed,
    required this.comment,
  });

  final double? points;
  final bool? passed;
  final String comment;

  Map<String, Object?> toJson() => {
    'points': points,
    'passed': passed,
    'comment': comment,
  };
}
