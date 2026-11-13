class SolutionReview {
  const SolutionReview({
    required this.points,
    required this.passed,
    required this.comment,
  });

  final int points;
  final List<String> passed;
  final String comment;
}
