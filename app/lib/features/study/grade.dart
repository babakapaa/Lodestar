import 'grade_parts.dart';

class Grade {
  const Grade({
    required this.studentId,
    required this.points,
    required this.passed,
    required this.comment,
    required this.gradedAt,
    this.parts,
  });

  factory Grade.fromJson(Map<String, Object?> json) {
    return Grade(
      studentId: json['studentId']! as String,
      points: (json['points'] as num?)?.toDouble(),
      passed: json['passed'] as bool?,
      parts: GradeParts.maybe(json['parts']),
      comment: json['comment'] as String? ?? '',
      gradedAt: DateTime.parse(json['gradedAt']! as String),
    );
  }

  static Grade? maybe(Object? json) =>
      json is Map<String, Object?> ? Grade.fromJson(json) : null;

  final String studentId;
  final double? points;
  final bool? passed;
  final GradeParts? parts;
  final String comment;
  final DateTime gradedAt;
}
