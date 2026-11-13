import 'defense_question.dart';
import 'defense_status.dart';

class Defense {
  const Defense({
    required this.studentId,
    required this.selected,
    required this.status,
    required this.questions,
    required this.summary,
    required this.heldAt,
  });

  factory Defense.fromJson(Map<String, Object?> json) {
    final heldAt = json['heldAt'];
    return Defense(
      studentId: json['studentId']! as String,
      selected: json['selected'] as bool? ?? false,
      status: DefenseStatus.parse(json['status']),
      questions: [
        for (final item in json['questions'] as List<Object?>? ?? const [])
          DefenseQuestion.fromJson(item! as Map<String, Object?>),
      ],
      summary: json['summary'] as String? ?? '',
      heldAt: heldAt is String ? DateTime.parse(heldAt) : null,
    );
  }

  final String studentId;
  final bool selected;
  final DefenseStatus status;
  final List<DefenseQuestion> questions;
  final String summary;
  final DateTime? heldAt;
}
