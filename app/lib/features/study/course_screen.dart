import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../attendance/attendance_providers.dart';
import 'assignments_section.dart';
import 'loaded_list_view.dart';
import 'navigator_tile.dart';
import 'study_providers.dart';

class CourseScreen extends ConsumerWidget {
  const CourseScreen({required this.courseId, super.key});

  final String courseId;

  Future<void> _refresh(WidgetRef ref) {
    ref.invalidate(assignmentsProvider(courseId));
    return ref.refresh(courseProvider(courseId).future);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final details = ref.watch(courseProvider(courseId));
    final attendance = ref.watch(attendanceEnabledProvider);
    final course = details.value?.value;
    final canCreate =
        course != null && course.course.role.teaches && !course.course.archived;
    return Scaffold(
      appBar: AppBar(title: Text(course?.course.title ?? '')),
      floatingActionButton: canCreate
          ? FloatingActionButton.extended(
              onPressed: () => context.push(AppRoutes.assignmentNew(courseId)),
              icon: const Icon(Icons.add),
              label: Text(l10n.assignmentNew),
            )
          : null,
      body: SafeArea(
        top: false,
        child: LoadedListView(
          value: details,
          bottomInset: canCreate ? 80 : 0,
          onRefresh: () => _refresh(ref),
          builder: (details) => [
            if (details.course.description.isNotEmpty) ...[
              SelectableText(details.course.description),
              const SizedBox(height: Spacing.xl),
            ],
            if (details.teachers.isNotEmpty) ...[
              Section(
                title: l10n.courseTeachers,
                child: Text(
                  [for (final person in details.teachers) person.label]
                      .join(', '),
                ),
              ),
              const SizedBox(height: Spacing.xl),
            ],
            if (details.units.isNotEmpty) ...[
              Section(
                title: l10n.courseGroups,
                child: Text(details.units.join(', ')),
              ),
              const SizedBox(height: Spacing.xl),
            ],
            NavigatorTile(
              icon: Icons.explore_outlined,
              title: l10n.navigatorAsk,
              route: AppRoutes.newDialog(
                scopeKind: 'course',
                scopeId: courseId,
              ),
            ),
            const SizedBox(height: Spacing.m),
            NavigatorTile(
              icon: Icons.forum_outlined,
              title: l10n.chatsCourseTitle,
              route: AppRoutes.courseChats(
                courseId,
                teach: details.course.role.teaches,
              ),
            ),
            const SizedBox(height: Spacing.m),
            NavigatorTile(
              key: const ValueKey('course-league'),
              icon: Icons.leaderboard_outlined,
              title: l10n.leagueTitle,
              subtitle: l10n.leagueHint,
              route: AppRoutes.league('study.course', courseId),
            ),
            if (details.course.role.teaches) ...[
              const SizedBox(height: Spacing.m),
              NavigatorTile(
                icon: Icons.library_books_outlined,
                title: l10n.materialsTitle,
                subtitle: l10n.materialsHint,
                route: AppRoutes.courseMaterials(courseId),
              ),
              const SizedBox(height: Spacing.m),
              NavigatorTile(
                key: const ValueKey('course-lms'),
                icon: Icons.table_chart_outlined,
                title: l10n.courseTablesTitle,
                subtitle: l10n.courseTablesHint,
                route: AppRoutes.courseLms(courseId),
              ),
              if (attendance) ...[
                const SizedBox(height: Spacing.m),
                NavigatorTile(
                  key: const ValueKey('course-attendance'),
                  icon: Icons.how_to_reg_outlined,
                  title: l10n.attendanceTitle,
                  subtitle: l10n.attendanceHint,
                  route: AppRoutes.attendanceCourse(courseId),
                ),
              ],
            ],
            const SizedBox(height: Spacing.xl),
            AssignmentsSection(
              courseId: courseId,
              teaches: details.course.role.teaches,
            ),
          ],
        ),
      ),
    );
  }
}
