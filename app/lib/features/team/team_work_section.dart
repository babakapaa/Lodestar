import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../profile/profile_providers.dart';
import '../study/assignment_details.dart';
import '../study/navigator_tile.dart';
import 'proposal_status.dart';
import 'team_proposal_card.dart';
import 'team_propose_form.dart';
import 'team_providers.dart';

class TeamWorkSection extends ConsumerWidget {
  const TeamWorkSection({required this.details, super.key});

  static const scopeKind = 'study.assignment';

  final AssignmentDetails details;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final id = details.assignment.id;
    final scope = ref.watch(teamScopeProvider((scopeKind, id))).value;
    final myTeam = scope?.myTeam;
    if (scope == null) {
      return const SizedBox.shrink();
    }
    if (myTeam == null) {
      return NavigatorTile(
        key: const ValueKey('team-find'),
        icon: Icons.groups_outlined,
        title: scope.solo ? l10n.teamSoloTitle : l10n.teamFind,
        subtitle: scope.solo ? l10n.teamSolo : l10n.teamFindHint,
        route: AppRoutes.teamScope(scopeKind, id),
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
              (proposal == null || proposal.status != ProposalStatus.sent)) ...[
            const SizedBox(height: Spacing.l),
            TeamProposeForm(assignment: details, team: team),
          ],
        ],
      ),
    );
  }
}
