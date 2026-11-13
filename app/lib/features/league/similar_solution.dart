class SimilarSolution {
  const SimilarSolution({
    required this.taskTitle,
    required this.teamName,
    required this.text,
  });

  factory SimilarSolution.fromJson(Map<String, Object?> json) =>
      SimilarSolution(
        taskTitle: json['taskTitle'] as String? ?? '',
        teamName: json['teamName'] as String? ?? '',
        text: json['text'] as String? ?? '',
      );

  final String taskTitle;
  final String teamName;
  final String text;
}
