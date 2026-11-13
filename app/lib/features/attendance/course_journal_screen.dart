import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/async_list_view.dart';
import '../../design/components/notice_banner.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'attendance_providers.dart';
import 'course_journal.dart';
import 'journal_export_menu.dart';
import 'journal_lesson.dart';
import 'journal_lesson_tile.dart';
import 'journal_student_tile.dart';

class CourseJournalScreen extends ConsumerWidget {
  const CourseJournalScreen({required this.courseId, super.key});

  final String courseId;

  Future<void> _open(
    BuildContext context,
    WidgetRef ref,
    JournalLesson lesson,
  ) async {
    var id = lesson.sessionId ?? '';
    if (id.isEmpty) {
      await TeamAction.run(context, () async {
        id =
            (await ref
                    .read(attendanceApiProvider)
                    .open(lesson.seriesId, lesson.date))
                .id;
      });
    }
    if (id.isEmpty || !context.mounted) {
      return;
    }
    await context.push(AppRoutes.attendanceSession(id));
    ref.invalidate(courseJournalProvider(courseId));
  }

  List<Widget> _body(BuildContext context, WidgetRef ref, CourseJournal j) {
    final l10n = context.l10n;
    final away = j.longAbsent;
    return [
      if (away.isNotEmpty) ...[
        NoticeBanner(
          icon: Icons.person_off_outlined,
          message: l10n.attendanceLongAbsentHint(j.absentAlert),
        ),
        Section(
          title: l10n.attendanceLongAbsent,
          child: Column(
            children: [
              for (final s in away)
                JournalStudentTile(
                  key: ValueKey('away-${s.userId}'),
                  student: s,
                ),
            ],
          ),
        ),
        const SizedBox(height: Spacing.xl),
      ],
      Section(
        title: l10n.attendanceLessons,
        child: j.lessons.isEmpty
            ? Text(l10n.attendanceNoLessons)
            : Column(
                children: [
                  for (final lesson in j.lessons.reversed)
                    JournalLessonTile(
                      key: ValueKey('lesson-${lesson.seriesId}-${lesson.date}'),
                      lesson: lesson,
                      onOpen: () => _open(context, ref, lesson),
                    ),
                ],
              ),
      ),
      const SizedBox(height: Spacing.xl),
      Section(
        title: l10n.attendanceStudents,
        child: j.students.isEmpty
            ? Text(l10n.attendanceNoStudents)
            : Column(
                children: [
                  for (final s in j.students)
                    JournalStudentTile(key: ValueKey(s.userId), student: s),
                ],
              ),
      ),
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final journal = ref.watch(courseJournalProvider(courseId));
    return Scaffold(
      appBar: AppBar(
        title: Text(journal.value?.courseTitle ?? context.l10n.attendanceTitle),
        actions: [JournalExportMenu(courseId: courseId)],
      ),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: journal,
          onRefresh: () => ref.refresh(courseJournalProvider(courseId).future),
          builder: (j) => _body(context, ref, j),
        ),
      ),
    );
  }
}
