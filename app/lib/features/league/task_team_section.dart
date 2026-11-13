import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../profile/profile_providers.dart';
import '../study/navigator_tile.dart';
import '../team/proposal_status.dart';
import '../team/team_proposal_card.dart';
import '../team/team_providers.dart';
import 'league_solution_form.dart';
import 'league_task.dart';

class TaskTeamSection extends ConsumerWidget {
  const TaskTeamSection({required this.task, super.key});

  static const scopeKind = 'league.task';

  final LeagueTask task;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final scope = ref.watch(teamScopeProvider((scopeKind, task.id))).value;
    if (scope == null) {
      return const SizedBox.shrink();
    }
    final myTeam = scope.myTeam;
    if (myTeam == null) {
      return NavigatorTile(
        key: const ValueKey('team-find'),
        icon: Icons.groups_outlined,
        title: l10n.teamFind,
        subtitle: l10n.leagueTaskTeamHint,
        route: AppRoutes.teamScope(scopeKind, task.id),
      );
    }
    final team = ref.watch(teamProvider(myTeam.id)).value;
    final me = ref.watch(accountProvider).value?.id;
    final proposal = team?.proposal;
    return Section(
      title: l10n.teamWork,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          NavigatorTile(
            key: const ValueKey('my-team'),
            icon: Icons.groups_outlined,
            title: myTeam.name,
            subtitle: myTeam.names,
            route: AppRoutes.team(myTeam.id),
          ),
          if (team != null && proposal != null) ...[
            const SizedBox(height: Spacing.l),
            TeamProposalCard(details: team, proposal: proposal, me: me),
          ],
          if (team != null &&
              !task.closed &&
              (proposal == null || proposal.status != ProposalStatus.sent)) ...[
            const SizedBox(height: Spacing.l),
            LeagueSolutionForm(team: team),
          ],
        ],
      ),
    );
  }
}
