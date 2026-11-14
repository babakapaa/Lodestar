import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../core/navigation/tab_page.dart';
import '../../design/components/tile_group.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../admin/admin_providers.dart';
import 'course_tile.dart';
import 'loaded_list_view.dart';
import 'navigator_tile.dart';
import 'outbox_section.dart';
import 'study_providers.dart';

class StudyScreen extends ConsumerWidget {
  const StudyScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final sandbox =
        ref
            .watch(modulesProvider)
            .value
            ?.any((m) => m.name == 'sandbox' && m.enabled) ??
        false;
    return TabPage(
      title: l10n.tabStudy,
      child: LoadedListView(
        value: ref.watch(myCoursesProvider),
        onRefresh: () => ref.refresh(myCoursesProvider.future),
        builder: (courses) => [
          const OutboxSection(),
          if (courses.isEmpty)
            EmptyState(
              icon: Icons.school_outlined,
              title: l10n.studyNoCourses,
              body: l10n.studyNoCoursesBody,
            )
          else
            Section(
              title: l10n.studyCourses,
              child: TileGroup(
                children: [
                  for (final course in courses)
                    CourseTile(key: ValueKey(course.id), course: course),
                ],
              ),
            ),
          const SizedBox(height: Spacing.xl),
          NavigatorTile(
            key: const ValueKey('my-teams'),
            icon: Icons.groups_outlined,
            title: l10n.teamsTitle,
            subtitle: l10n.teamsHint,
            route: AppRoutes.teams,
          ),
          if (sandbox)
            NavigatorTile(
              key: const ValueKey('my-sandbox'),
              icon: Icons.memory,
              title: l10n.sandboxTitle,
              subtitle: l10n.sandboxHint,
              route: AppRoutes.sandbox,
            ),
        ],
      ),
    );
  }
}
