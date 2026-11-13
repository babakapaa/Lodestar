import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/study_labels_context.dart';
import 'attendance_labels_context.dart';
import 'attendance_providers.dart';

class OpenLessonsSection extends ConsumerWidget {
  const OpenLessonsSection({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    if (!ref.watch(attendanceEnabledProvider)) {
      return const SizedBox.shrink();
    }
    final lessons = ref.watch(openLessonsProvider).value ?? const [];
    if (lessons.isEmpty) {
      return const SizedBox.shrink();
    }
    final l10n = context.l10n;
    final labels = context.studyLabels;
    return Padding(
      padding: const EdgeInsets.only(bottom: Spacing.xl),
      child: Section(
        title: l10n.attendanceOpenSection,
        child: Column(
          children: [
            for (final lesson in lessons)
              ListTile(
                key: ValueKey('open-${lesson.sessionId}'),
                contentPadding: EdgeInsets.zero,
                leading: Icon(
                  lesson.status == null
                      ? Icons.how_to_reg_outlined
                      : Icons.check_circle_outline,
                ),
                title: Text(lesson.courseTitle),
                subtitle: Text(
                  [
                    '${labels.clock(lesson.startsAt)}–${labels.clock(lesson.endsAt)}',
                    if (lesson.room.isNotEmpty) lesson.room,
                    context.attendanceLabels.status(lesson.status),
                  ].join(', '),
                ),
                trailing: lesson.status == null
                    ? FilledButton(
                        onPressed: () => context.push(AppRoutes.attend),
                        child: Text(l10n.attendanceCheckIn),
                      )
                    : null,
              ),
          ],
        ),
      ),
    );
  }
}
