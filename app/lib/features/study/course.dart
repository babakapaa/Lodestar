import 'course_role.dart';

class Course {
  const Course({
    required this.id,
    required this.title,
    required this.code,
    required this.description,
    required this.role,
    required this.archived,
  });

  factory Course.fromJson(Map<String, Object?> json) {
    return Course(
      id: json['id']! as String,
      title: json['title']! as String,
      code: json['code'] as String?,
      description: json['description'] as String? ?? '',
      role: CourseRole.parse(json['role']),
      archived: json['archivedAt'] != null,
    );
  }

  final String id;
  final String title;
  final String? code;
  final String description;
  final CourseRole role;
  final bool archived;
}
