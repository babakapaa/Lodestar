class SetupProblem {
  const SetupProblem({
    required this.section,
    required this.code,
    required this.count,
  });

  factory SetupProblem.fromJson(Map<String, Object?> json) => SetupProblem(
    section: json['section']! as String,
    code: json['code']! as String,
    count: json['count'] as int? ?? 1,
  );

  final String section;
  final String code;
  final int count;
}
