import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import '../study/study_labels_context.dart';
import 'journal_lesson.dart';

class JournalLessonTile extends StatelessWidget {
  const JournalLessonTile({
    required this.lesson,
    required this.onOpen,
    super.key,
  });

  final JournalLesson lesson;
  final VoidCallback onOpen;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final state = lesson.cancelled
        ? l10n.attendanceLessonCancelled
        : lesson.sessionId == null
        ? l10n.attendanceLessonNone
        : lesson.open
        ? l10n.attendanceLessonOpen
        : l10n.attendanceHere(
            lesson.counts.here,
            lesson.counts.here +
                lesson.counts.absent +
                lesson.counts.excused +
                lesson.counts.unmarked,
          );
    return ListTile(
      contentPadding: EdgeInsets.zero,
      title: Text(
        '${labels.day(lesson.startsAt)}, ${labels.clock(lesson.startsAt)}',
      ),
      subtitle: Text(
        [if (lesson.room.isNotEmpty) lesson.room, state].join(', '),
      ),
      trailing: lesson.cancelled ? null : const Icon(Icons.chevron_right),
      onTap: lesson.cancelled ? null : onOpen,
    );
  }
}
