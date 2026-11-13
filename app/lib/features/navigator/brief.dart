class Brief {
  const Brief({
    required this.kind,
    required this.title,
    required this.course,
    required this.project,
  });

  factory Brief.fromJson(Map<String, Object?> json) => Brief(
    kind: json['kind']! as String,
    title: json['title']! as String,
    course: json['course'] as String? ?? '',
    project: json['mode'] == 'project',
  );

  static Brief? maybe(Object? json) =>
      json is Map<String, Object?> ? Brief.fromJson(json) : null;

  final String kind;
  final String title;
  final String course;
  final bool project;

  bool get isCourse => kind == 'course';
}
