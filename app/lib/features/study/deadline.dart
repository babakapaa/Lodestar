class Deadline {
  const Deadline({
    required this.assignmentId,
    required this.courseTitle,
    required this.title,
    required this.dueAt,
    required this.submittedAt,
  });

  factory Deadline.fromJson(Map<String, Object?> json) {
    final submittedAt = json['submittedAt'];
    return Deadline(
      assignmentId: json['assignmentId']! as String,
      courseTitle: json['courseTitle'] as String? ?? '',
      title: json['title']! as String,
      dueAt: DateTime.parse(json['dueAt']! as String),
      submittedAt: submittedAt is String ? DateTime.parse(submittedAt) : null,
    );
  }

  final String assignmentId;
  final String courseTitle;
  final String title;
  final DateTime dueAt;
  final DateTime? submittedAt;
}
