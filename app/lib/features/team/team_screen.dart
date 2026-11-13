import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../profile/profile_providers.dart';
import 'team_action.dart';
import 'team_details.dart';
import 'team_history_section.dart';
import 'team_invite_section.dart';
import 'team_leave_section.dart';
import 'team_members_section.dart';
import 'team_proposal_card.dart';
import 'team_providers.dart';
import 'team_refresh.dart';
import 'team_scope_header.dart';
import 'team_text_dialog.dart';

class TeamScreen extends ConsumerWidget {
  const TeamScreen({required this.teamId, super.key});

  final String teamId;

  Future<void> _rename(
    BuildContext context,
    WidgetRef ref,
    TeamDetails details,
  ) async {
    final l10n = context.l10n;
    final name = await TeamTextDialog.show(
      context,
      title: l10n.teamRename,
      label: l10n.teamName,
      action: l10n.actionSave,
      initial: details.team.name,
    );
    if (name == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).rename(details.team.id, name),
    );
    ref.refreshTeam(details);
  }

  List<Widget> _body(TeamDetails details, String? me) {
    final member = me != null && details.team.has(me);
    final proposal = details.proposal;
    return [
      TeamScopeHeader(scope: details.scope),
      const SizedBox(height: Spacing.xl),
      TeamMembersSection(details: details, me: me),
      if (proposal != null) ...[
        const SizedBox(height: Spacing.xl),
        TeamProposalCard(details: details, proposal: proposal, me: me),
      ],
      if (member || details.manage) ...[
        const SizedBox(height: Spacing.xl),
        TeamInviteSection(details: details),
      ],
      if (member) ...[
        const SizedBox(height: Spacing.xl),
        TeamLeaveSection(details: details),
      ],
      if (details.history.isNotEmpty) ...[
        const SizedBox(height: Spacing.xl),
        TeamHistorySection(history: details.history),
      ],
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final team = ref.watch(teamProvider(teamId));
    final me = ref.watch(accountProvider).value?.id;
    final loaded = team.value;
    final editable =
        loaded != null &&
        !loaded.team.archived &&
        (loaded.manage || (me != null && loaded.team.has(me)));
    return Scaffold(
      appBar: AppBar(
        title: Text(loaded?.team.name ?? ''),
        actions: [
          if (editable)
            IconButton(
              key: const ValueKey('team-rename'),
              tooltip: context.l10n.teamRename,
              icon: const Icon(Icons.edit_outlined),
              onPressed: () => _rename(context, ref, loaded),
            ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: team,
          onRefresh: () => ref.refresh(teamProvider(teamId).future),
          builder: (details) => _body(details, me),
        ),
      ),
    );
  }
}
