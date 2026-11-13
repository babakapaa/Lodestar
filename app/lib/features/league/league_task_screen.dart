import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/async_list_view.dart';
import '../../design/components/section.dart';
import '../../design/components/spaced_column.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../study/navigator_tile.dart';
import '../team/team_action.dart';
import '../team/team_text_dialog.dart';
import 'league_labels_context.dart';
import 'league_providers.dart';
import 'league_solution.dart';
import 'review_dialog.dart';
import 'solution_card.dart';
import 'task_details.dart';
import 'task_team_section.dart';

class LeagueTaskScreen extends ConsumerWidget {
  const LeagueTaskScreen({required this.taskId, super.key});

  final String taskId;

  Future<void> _review(
    BuildContext context,
    WidgetRef ref,
    TaskDetails details,
    LeagueSolution solution,
  ) async {
    final review = await ReviewDialog.show(context, details, solution);
    if (review == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref
          .read(leagueApiProvider)
          .review(
            solution.id,
            points: review.points,
            passed: review.passed,
            comment: review.comment,
          ),
    );
    ref.invalidate(leagueTaskProvider(taskId));
  }

  Future<void> _explain(
    BuildContext context,
    WidgetRef ref,
    LeagueSolution solution,
  ) async {
    final l10n = context.l10n;
    final text = await TeamTextDialog.show(
      context,
      title: l10n.leagueExplain,
      label: l10n.leagueExplainText,
      action: l10n.actionSave,
      initial: solution.explanation,
      lines: 3,
    );
    if (text == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(leagueApiProvider).explain(solution.id, text),
    );
    ref.invalidate(leagueTaskProvider(taskId));
  }

  List<Widget> _body(BuildContext context, WidgetRef ref, TaskDetails details) {
    final l10n = context.l10n;
    final labels = context.leagueLabels;
    final task = details.task;
    final textTheme = Theme.of(context).textTheme;
    final muted = textTheme.bodyMedium!.copyWith(color: context.palette.muted);
    return [
      Text(details.league, style: muted),
      Text(task.title, style: textTheme.headlineSmall),
      const SizedBox(height: Spacing.s),
      Text(
        [
          l10n.leaguePlus(task.points),
          if (task.partner.isNotEmpty) l10n.leagueTaskFrom(task.partner),
          task.minSize == task.maxSize
              ? l10n.teamSizeExact(task.minSize)
              : l10n.teamSizeRange(task.minSize, task.maxSize),
          if (task.closed) l10n.leagueTaskClosed else labels.due(task.dueAt),
        ].join(', '),
        style: muted,
      ),
      if (task.about.isNotEmpty) ...[
        const SizedBox(height: Spacing.l),
        SelectableText(task.about),
      ],
      const SizedBox(height: Spacing.xl),
      SpacedColumn(
        children: [
          if (!details.leads) TaskTeamSection(task: task),
          if (details.leads)
            NavigatorTile(
              key: const ValueKey('teams'),
              icon: Icons.groups_outlined,
              title: l10n.teamsTitle,
              route: AppRoutes.teamScope(TaskTeamSection.scopeKind, task.id),
            ),
          Section(
            title: l10n.leagueSolutions(details.solutions.length),
            child: SpacedColumn(
              gap: Spacing.m,
              children: [
                if (details.solutions.isEmpty) Text(l10n.leagueSolutionsNone),
                for (final solution in details.solutions)
                  SolutionCard(
                    details: details,
                    solution: solution,
                    onReview: details.leads
                        ? () => _review(context, ref, details, solution)
                        : null,
                    onExplain: details.leads
                        ? null
                        : () => _explain(context, ref, solution),
                  ),
              ],
            ),
          ),
        ],
      ),
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final details = ref.watch(leagueTaskProvider(taskId));
    final loaded = details.value;
    return Scaffold(
      appBar: AppBar(
        title: Text(context.l10n.leagueTitle),
        actions: [
          if (loaded != null && loaded.leads)
            IconButton(
              key: const ValueKey('task-edit'),
              tooltip: context.l10n.leagueTaskEdit,
              icon: const Icon(Icons.edit_outlined),
              onPressed: () => context.push(AppRoutes.leagueTaskEdit(taskId)),
            ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: details,
          onRefresh: () => ref.refresh(leagueTaskProvider(taskId).future),
          builder: (details) => _body(context, ref, details),
        ),
      ),
    );
  }
}
