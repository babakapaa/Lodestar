import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'lesson.dart';
import 'lesson_action.dart';
import 'lesson_status.dart';

class LessonActionsSheet extends StatelessWidget {
  const LessonActionsSheet({required this.lesson, super.key});

  final Lesson lesson;

  static Future<LessonAction?> show(BuildContext context, Lesson lesson) =>
      showModalBottomSheet<LessonAction>(
        context: context,
        builder: (_) => LessonActionsSheet(lesson: lesson),
      );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final planned = lesson.status == LessonStatus.planned;
    return SafeArea(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (lesson.status != LessonStatus.cancelled)
            ListTile(
              key: const ValueKey('lesson-cancel'),
              leading: const Icon(Icons.event_busy_outlined),
              title: Text(l10n.lessonCancel),
              onTap: () => Navigator.of(context).pop(LessonAction.cancel),
            ),
          ListTile(
            key: const ValueKey('lesson-move'),
            leading: const Icon(Icons.schedule_outlined),
            title: Text(l10n.lessonMove),
            onTap: () => Navigator.of(context).pop(LessonAction.move),
          ),
          ListTile(
            key: const ValueKey('lesson-room'),
            leading: const Icon(Icons.meeting_room_outlined),
            title: Text(l10n.lessonRoom),
            onTap: () => Navigator.of(context).pop(LessonAction.room),
          ),
          if (!planned)
            ListTile(
              key: const ValueKey('lesson-revert'),
              leading: const Icon(Icons.undo),
              title: Text(l10n.lessonRevert),
              onTap: () => Navigator.of(context).pop(LessonAction.revert),
            ),
        ],
      ),
    );
  }
}
