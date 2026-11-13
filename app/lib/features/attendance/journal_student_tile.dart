import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'journal_student.dart';

class JournalStudentTile extends StatelessWidget {
  const JournalStudentTile({required this.student, super.key});

  final JournalStudent student;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final rate = student.rate;
    return ListTile(
      contentPadding: EdgeInsets.zero,
      title: Text(student.name),
      subtitle: Text(
        [
          if (rate != null) l10n.attendanceRate((rate * 100).round()),
          if (student.streak > 0) l10n.attendanceStreak(student.streak),
        ].join(', '),
        style: student.longAbsent
            ? TextStyle(color: context.palette.warn)
            : null,
      ),
    );
  }
}
