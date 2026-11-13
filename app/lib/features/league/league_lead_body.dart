import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/spaced_column.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'award_grant.dart';
import 'dispute_answer_dialog.dart';
import 'league_dispute.dart';
import 'league_disputes_section.dart';
import 'league_grants_section.dart';
import 'league_moves_section.dart';
import 'league_players_section.dart';
import 'league_providers.dart';
import 'league_refresh.dart';
import 'league_tasks_section.dart';
import 'league_transition.dart';
import 'league_view.dart';

class LeagueLeadBody extends ConsumerWidget {
  const LeagueLeadBody({required this.view, super.key});

  final LeagueView view;

  Future<void> _run(
    BuildContext context,
    WidgetRef ref,
    Future<void> Function() action,
  ) async {
    await TeamAction.run(context, action);
    ref.refreshLeague(view.kind, view.originId);
  }

  Future<void> _closeSeason(BuildContext context, WidgetRef ref) async {
    final l10n = context.l10n;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.leagueCloseSeason),
        content: Text(l10n.leagueCloseSeasonHint),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: Text(l10n.actionCancel),
          ),
          FilledButton(
            key: const ValueKey('season-close-confirm'),
            onPressed: () => Navigator.of(context).pop(true),
            child: Text(l10n.leagueCloseSeason),
          ),
        ],
      ),
    );
    if ((confirmed ?? false) && context.mounted) {
      await _run(
        context,
        ref,
        () => ref.read(leagueApiProvider).closeSeason(view.id),
      );
    }
  }

  Future<void> _answer(
    BuildContext context,
    WidgetRef ref,
    LeagueDispute dispute,
  ) async {
    final answer = await DisputeAnswerDialog.show(context, dispute);
    if (answer == null || !context.mounted) {
      return;
    }
    await _run(
      context,
      ref,
      () => ref
          .read(leagueApiProvider)
          .answer(
            dispute.id,
            accept: answer.accept,
            answer: answer.answer,
            points: answer.points,
          ),
    );
  }

  Widget _decision(BuildContext context, WidgetRef ref, LeagueTransition move) {
    final l10n = context.l10n;
    final api = ref.read(leagueApiProvider);
    return Wrap(
      children: [
        IconButton(
          key: ValueKey('move-approve-${move.id}'),
          tooltip: l10n.leagueAccept,
          icon: const Icon(Icons.check),
          onPressed: () =>
              _run(context, ref, () => api.decide(move.id, approve: true)),
        ),
        IconButton(
          tooltip: l10n.leagueReject,
          icon: const Icon(Icons.close),
          onPressed: () =>
              _run(context, ref, () => api.decide(move.id, approve: false)),
        ),
      ],
    );
  }

  String _name(String userId) {
    for (final player in view.players) {
      if (player.person.userId == userId) {
        return player.person.label;
      }
    }
    return userId;
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    return SpacedColumn(
      children: [
        Wrap(
          spacing: Spacing.s,
          runSpacing: Spacing.s,
          children: [
            LodestarButton(
              key: const ValueKey('league-rules'),
              label: l10n.leagueRules,
              icon: Icons.tune,
              variant: LodestarButtonVariant.secondary,
              onPressed: () =>
                  router.push(AppRoutes.leagueRules(view.kind, view.originId)),
            ),
            LodestarButton(
              key: const ValueKey('season-close'),
              label: l10n.leagueCloseSeason,
              icon: Icons.flag_circle_outlined,
              variant: LodestarButtonVariant.secondary,
              onPressed: () => _closeSeason(context, ref),
            ),
          ],
        ),
        if (view.pending.isNotEmpty)
          LeagueMovesSection(
            view: view,
            moves: view.pending,
            title: l10n.leaguePendingMoves,
            trailing: (move) => _decision(context, ref, move),
          ),
        if (view.openDisputes.isNotEmpty)
          LeagueDisputesSection(
            disputes: view.openDisputes,
            onAnswer: (dispute) => _answer(context, ref, dispute),
          ),
        LeagueTasksSection(
          tasks: view.tasks,
          action: Align(
            alignment: Alignment.centerLeft,
            child: TextButton.icon(
              key: const ValueKey('task-new'),
              onPressed: () => router.push(
                AppRoutes.leagueTaskNew(view.kind, view.originId),
              ),
              icon: const Icon(Icons.add),
              label: Text(l10n.leagueTaskNew),
            ),
          ),
        ),
        LeaguePlayersSection(view: view),
        if (view.grants.isNotEmpty)
          LeagueGrantsSection(
            grants: view.grants,
            nameOf: _name,
            onRevoke: (AwardGrant grant) => _run(
              context,
              ref,
              () => ref.read(leagueApiProvider).revoke(grant.id),
            ),
          ),
      ],
    );
  }
}
