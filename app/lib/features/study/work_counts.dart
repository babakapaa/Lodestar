class WorkCounts {
  const WorkCounts({
    required this.students,
    required this.submitted,
    required this.graded,
  });

  factory WorkCounts.fromJson(Map<String, Object?> json) {
    return WorkCounts(
      students: (json['students']! as num).toInt(),
      submitted: (json['submitted']! as num).toInt(),
      graded: (json['graded']! as num).toInt(),
    );
  }

  static WorkCounts? maybe(Object? json) =>
      json is Map<String, Object?> ? WorkCounts.fromJson(json) : null;

  final int students;
  final int submitted;
  final int graded;
}
