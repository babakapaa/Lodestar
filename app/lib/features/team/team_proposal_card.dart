import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'proposal_status.dart';
import 'team_action.dart';
import 'team_details.dart';
import 'team_labels_context.dart';
import 'team_proposal.dart';
import 'team_providers.dart';
import 'team_refresh.dart';
import 'team_text_dialog.dart';
import 'vote_decision.dart';

class TeamProposalCard extends ConsumerWidget {
  const TeamProposalCard({
    required this.details,
    required this.proposal,
    required this.me,
    super.key,
  });

  final TeamDetails details;
  final TeamProposal proposal;
  final String? me;

  Future<void> _vote(
    BuildContext context,
    WidgetRef ref, {
    required bool agree,
  }) async {
    final l10n = context.l10n;
    var comment = '';
    if (!agree) {
      final text = await TeamTextDialog.show(
        context,
        title: l10n.teamObject,
        label: l10n.teamObjectReason,
        action: l10n.teamObject,
        lines: 2,
      );
      if (text == null || !context.mounted) {
        return;
      }
      comment = text;
    }
    await TeamAction.run(
      context,
      () => ref
          .read(teamApiProvider)
          .vote(proposal.id, agree: agree, comment: comment),
    );
    ref.refreshTeam(details);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final labels = context.teamLabels;
    final textTheme = Theme.of(context).textTheme;
    final muted = textTheme.bodyMedium!.copyWith(color: context.palette.muted);
    final mine = me == null ? null : proposal.decisionOf(me!);
    final open = proposal.status == ProposalStatus.open;
    return LodestarCard(
      key: const ValueKey('team-proposal'),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(
            open ? l10n.teamProposalOpen : l10n.teamProposalSent,
            style: textTheme.titleMedium,
          ),
          const SizedBox(height: Spacing.s),
          Text(
            [
              if (proposal.proposedBy case final by?)
                l10n.teamProposedBy(by.label),
              if (open)
                l10n.teamProposalSilentAt(labels.time(proposal.silentAt)),
            ].join('\n'),
            style: muted,
          ),
          if (proposal.note.isNotEmpty) ...[
            const SizedBox(height: Spacing.s),
            Text(proposal.note),
          ],
          if (proposal.work.text.isNotEmpty) ...[
            const SizedBox(height: Spacing.m),
            SelectableText(proposal.work.text),
          ],
          for (final link in proposal.work.links) SelectableText(link),
          const SizedBox(height: Spacing.m),
          for (final vote in proposal.votes)
            ListTile(
              dense: true,
              contentPadding: EdgeInsets.zero,
              leading: Icon(switch (vote.decision) {
                VoteDecision.agreed => Icons.check_circle_outline,
                VoteDecision.objected => Icons.block,
                VoteDecision.silent => Icons.schedule,
                VoteDecision.waiting => Icons.hourglass_empty,
              }),
              title: Text(vote.person.label),
              subtitle: Text(
                [
                  labels.decision(vote.decision),
                  if (vote.comment.isNotEmpty) vote.comment,
                ].join(': '),
              ),
            ),
          if (open && mine != null) ...[
            const SizedBox(height: Spacing.s),
            Wrap(
              spacing: Spacing.s,
              children: [
                FilledButton.icon(
                  key: const ValueKey('vote-agree'),
                  onPressed: mine == VoteDecision.agreed
                      ? null
                      : () => _vote(context, ref, agree: true),
                  icon: const Icon(Icons.check),
                  label: Text(l10n.teamAgree),
                ),
                OutlinedButton.icon(
                  key: const ValueKey('vote-object'),
                  onPressed: mine == VoteDecision.objected
                      ? null
                      : () => _vote(context, ref, agree: false),
                  icon: const Icon(Icons.block),
                  label: Text(l10n.teamObject),
                ),
              ],
            ),
          ],
        ],
      ),
    );
  }
}
