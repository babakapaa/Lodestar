class PilotCourse {
  const PilotCourse({required this.id, required this.title, this.code});

  factory PilotCourse.fromJson(Map<String, Object?> json) => PilotCourse(
    id: json['id']! as String,
    title: json['title']! as String,
    code: json['code'] as String?,
  );

  final String id;
  final String title;
  final String? code;
}
