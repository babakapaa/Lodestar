import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/tokens_context.dart';
import 'lesson.dart';
import 'lesson_status.dart';
import 'study_labels_context.dart';

class LessonTile extends StatelessWidget {
  const LessonTile({required this.lesson, this.onTap, super.key});

  final Lesson lesson;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final labels = context.studyLabels;
    final palette = context.palette;
    final textTheme = Theme.of(context).textTheme;
    final cancelled = lesson.status == LessonStatus.cancelled;
    final status = labels.status(lesson.status);
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: SizedBox(
        width: 52,
        child: Text(labels.clock(lesson.startsAt), style: textTheme.titleSmall),
      ),
      title: Text(
        lesson.courseTitle,
        style: cancelled
            ? const TextStyle(decoration: TextDecoration.lineThrough)
            : null,
      ),
      subtitle: Text(
        [
          '${labels.clock(lesson.startsAt)}–${labels.clock(lesson.endsAt)}',
          labels.kind(lesson.kind),
          if (lesson.room.isNotEmpty) lesson.room,
          ?status,
          if (lesson.note.isNotEmpty) lesson.note,
        ].join(', '),
        style: status == null ? null : TextStyle(color: palette.warn),
      ),
      onTap: onTap ?? () => context.push(AppRoutes.course(lesson.courseId)),
    );
  }
}
