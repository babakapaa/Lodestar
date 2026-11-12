import 'assignment.dart';
import 'grade.dart';
import 'work_counts.dart';

class AssignmentItem {
  const AssignmentItem({
    required this.assignment,
    required this.submittedAt,
    required this.grade,
    required this.counts,
  });

  factory AssignmentItem.fromJson(Map<String, Object?> json) {
    final mine = json['mine'];
    final progress = mine is Map<String, Object?>
        ? mine
        : const <String, Object?>{};
    final submittedAt = progress['submittedAt'];
    return AssignmentItem(
      assignment: Assignment.fromJson(json),
      submittedAt: submittedAt is String ? DateTime.parse(submittedAt) : null,
      grade: Grade.maybe(progress['grade']),
      counts: WorkCounts.maybe(json['counts']),
    );
  }

  final Assignment assignment;
  final DateTime? submittedAt;
  final Grade? grade;
  final WorkCounts? counts;
}
