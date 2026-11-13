import '../study/grade_input.dart';
import 'defense_status.dart';

class FinishChoice {
  const FinishChoice({required this.status, required this.grade});

  final DefenseStatus status;
  final GradeInput? grade;
}
