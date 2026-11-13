import '../study/course_person.dart';
import '../study/grade.dart';
import 'defense_slot.dart';
import 'defense_status.dart';

class DefenseRow {
  const DefenseRow({
    required this.student,
    required this.selected,
    required this.status,
    required this.slot,
    required this.grade,
  });

  factory DefenseRow.fromJson(Map<String, Object?> json) {
    final slot = json['slot'];
    return DefenseRow(
      student: CoursePerson.fromJson(json['student']! as Map<String, Object?>),
      selected: json['selected'] as bool? ?? false,
      status: DefenseStatus.parse(json['status']),
      slot: slot is Map<String, Object?> ? DefenseSlot.fromJson(slot) : null,
      grade: Grade.maybe(json['grade']),
    );
  }

  final CoursePerson student;
  final bool selected;
  final DefenseStatus status;
  final DefenseSlot? slot;
  final Grade? grade;
}
