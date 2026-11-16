import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/layout_size.dart';
import '../../core/navigation/tab_page.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../design/components/tile_group.dart';
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
import '../study/study_providers.dart';
import 'today_actions.dart';
import 'today_greeting.dart';
import 'today_messages.dart';
import 'today_notes.dart';

class TodayScreen extends ConsumerWidget {
  const TodayScreen({super.key});

  static const sideWidth = 380.0;

  List<Widget> _agenda(BuildContext context, Agenda agenda) {
    final l10n = context.l10n;
    return [
      if (agenda.lessons.isEmpty && agenda.deadlines.isEmpty)
        Section(
          title: l10n.todayPlan,
          child: LodestarCard(
            child: EmptyState(
              icon: Icons.wb_sunny_outlined,
              title: l10n.emptyTodayTitle,
              body: l10n.emptyTodayBody,
            ),
          ),
        ),
      if (agenda.lessons.isNotEmpty) ...[
        Section(
          title: l10n.todayLessons,
          child: TileGroup(
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
          child: TileGroup(
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
    ];
  }

  List<Widget> _side() => const [
    TodayMessages(),
    SizedBox(height: Spacing.xl),
    TodayNotes(),
  ];

  List<Widget> _body(BuildContext context, Agenda agenda, bool wide) {
    return [
      const SetupBanner(),
      const OutboxSection(),
      const OpenLessonsSection(),
      TodayGreeting(date: agenda.date),
      const SizedBox(height: Spacing.xl),
      TodayActions(date: agenda.date),
      const SizedBox(height: Spacing.xxl),
      if (wide)
        Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: _agenda(context, agenda),
              ),
            ),
            const SizedBox(width: Spacing.xxl),
            SizedBox(
              width: sideWidth,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: _side(),
              ),
            ),
          ],
        )
      else ...[
        ..._agenda(context, agenda),
        ..._side(),
      ],
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final wide =
        LayoutSize.of(MediaQuery.sizeOf(context).width) == LayoutSize.expanded;
    return TabPage(
      title: context.l10n.tabToday,
      child: LoadedListView(
        value: ref.watch(agendaProvider(null)),
        maxWidth: wide ? 1160 : 720,
        onRefresh: () {
          ref.invalidate(openLessonsProvider);
          return ref.refresh(agendaProvider(null).future);
        },
        builder: (agenda) => _body(context, agenda, wide),
      ),
    );
  }
}
