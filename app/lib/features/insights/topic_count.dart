class TopicCount {
  const TopicCount({
    required this.materialId,
    required this.title,
    required this.students,
    required this.recent,
  });

  factory TopicCount.fromJson(Map<String, Object?> json) => TopicCount(
    materialId: json['materialId']! as String,
    title: json['title'] as String? ?? '',
    students: json['students'] as int? ?? 0,
    recent: json['recent'] as int? ?? 0,
  );

  final String materialId;
  final String title;
  final int students;
  final int recent;
}
