class SetupLegalAnswer {
  const SetupLegalAnswer({
    required this.question,
    required this.title,
    required this.answer,
    required this.basis,
    required this.approved,
  });

  factory SetupLegalAnswer.fromJson(Map<String, Object?> json) =>
      SetupLegalAnswer(
        question: json['question']! as String,
        title: json['title'] as String? ?? '',
        answer: json['answer'] as String? ?? '',
        basis: json['basis'] as String? ?? '',
        approved: json['approvedAt'] != null,
      );

  final String question;
  final String title;
  final String answer;
  final String basis;
  final bool approved;

  bool get answered => answer.isNotEmpty;
}
