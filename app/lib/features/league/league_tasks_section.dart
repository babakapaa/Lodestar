import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/navigator_tile.dart';
import 'league_labels_context.dart';
import 'league_task.dart';

class LeagueTasksSection extends StatelessWidget {
  const LeagueTasksSection({required this.tasks, this.action, super.key});

  final List<LeagueTask> tasks;
  final Widget? action;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.leagueLabels;
    final extra = action;
    return Section(
      title: l10n.leagueTasks,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          if (tasks.isEmpty) Text(l10n.leagueTasksNone),
          for (final task in tasks) ...[
            NavigatorTile(
              key: ValueKey('task-${task.id}'),
              icon: task.closed ? Icons.task_alt : Icons.flag_outlined,
              title: task.title,
              subtitle: [
                l10n.leaguePlus(task.points),
                if (task.partner.isNotEmpty) task.partner,
                if (task.closed)
                  l10n.leagueTaskClosed
                else
                  labels.due(task.dueAt),
              ].join(', '),
              route: AppRoutes.leagueTask(task.id),
            ),
            const SizedBox(height: Spacing.s),
          ],
          ?extra,
        ],
      ),
    );
  }
}
