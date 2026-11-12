import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'assignment.dart';
import 'student_work.dart';
import 'study_labels_context.dart';

class StudentWorkTile extends StatelessWidget {
  const StudentWorkTile({
    required this.work,
    required this.assignment,
    super.key,
  });

  final StudentWork work;
  final Assignment assignment;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final latest = work.latest;
    final grade = work.grade;
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Icon(latest == null ? Icons.hourglass_empty : Icons.task_alt),
      title: Text(work.student.label),
      subtitle: Text(
        latest == null
            ? l10n.workNotSubmitted
            : l10n.myWorkSentAt(labels.time(latest.submittedAt)),
      ),
      trailing: grade == null ? null : Text(labels.grade(grade, assignment)),
    );
  }
}
