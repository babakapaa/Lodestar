import 'course_person.dart';
import 'grade.dart';
import 'submission.dart';

class StudentWork {
  const StudentWork({
    required this.student,
    required this.latest,
    required this.attempts,
    required this.grade,
  });

  factory StudentWork.fromJson(Map<String, Object?> json) {
    return StudentWork(
      student: CoursePerson.fromJson(json['student']! as Map<String, Object?>),
      latest: Submission.maybe(json['latest']),
      attempts: (json['attempts'] as num? ?? 0).toInt(),
      grade: Grade.maybe(json['grade']),
    );
  }

  static StudentWork? maybe(Object? json) =>
      json is Map<String, Object?> ? StudentWork.fromJson(json) : null;

  final CoursePerson student;
  final Submission? latest;
  final int attempts;
  final Grade? grade;
}
