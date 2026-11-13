import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../core/navigation/tab_page.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../attendance/attendance_providers.dart';
import '../attendance/open_lessons_section.dart';
import '../setup/setup_banner.dart';
import '../study/agenda.dart';
import '../study/deadline_tile.dart';
import '../study/lesson_tile.dart';
import '../study/loaded_list_view.dart';
import '../study/outbox_section.dart';
import '../study/study_day.dart';
import '../study/study_labels_context.dart';
import '../study/study_providers.dart';

class TodayScreen extends ConsumerWidget {
  const TodayScreen({super.key});

  List<Widget> _body(BuildContext context, Agenda agenda) {
    final l10n = context.l10n;
    return [
      const SetupBanner(),
      const OutboxSection(),
      const OpenLessonsSection(),
      Text(
        context.studyLabels.day(StudyDay.parse(agenda.date)),
        style: Theme.of(context).textTheme.titleLarge,
      ),
      const SizedBox(height: Spacing.l),
      if (agenda.lessons.isEmpty && agenda.deadlines.isEmpty)
        EmptyState(
          icon: Icons.wb_sunny_outlined,
          title: l10n.emptyTodayTitle,
          body: l10n.emptyTodayBody,
        ),
      if (agenda.lessons.isNotEmpty) ...[
        Section(
          title: l10n.todayLessons,
          child: Column(
            children: [
              for (final lesson in agenda.lessons)
                LessonTile(
                  key: ValueKey('${lesson.seriesId}/${lesson.date}'),
                  lesson: lesson,
                ),
            ],
          ),
        ),
        const SizedBox(height: Spacing.xl),
      ],
      if (agenda.deadlines.isNotEmpty) ...[
        Section(
          title: l10n.todayDeadlines,
          child: Column(
            children: [
              for (final deadline in agenda.deadlines)
                DeadlineTile(
                  key: ValueKey(deadline.assignmentId),
                  deadline: deadline,
                ),
            ],
          ),
        ),
        const SizedBox(height: Spacing.xl),
      ],
      Align(
        alignment: Alignment.centerLeft,
        child: LodestarButton(
          label: l10n.todayWeek,
          icon: Icons.calendar_month_outlined,
          variant: LodestarButtonVariant.secondary,
          onPressed: () => context.push(AppRoutes.schedule(agenda.date)),
        ),
      ),
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return TabPage(
      title: context.l10n.tabToday,
      child: LoadedListView(
        value: ref.watch(agendaProvider(null)),
        onRefresh: () {
          ref.invalidate(openLessonsProvider);
          return ref.refresh(agendaProvider(null).future);
        },
        builder: (agenda) => _body(context, agenda),
      ),
    );
  }
}
