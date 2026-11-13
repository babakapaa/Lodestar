import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/spaced_column.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import '../team/team_text_dialog.dart';
import 'league_disputes_section.dart';
import 'league_display_dialog.dart';
import 'league_entries_section.dart';
import 'league_entry.dart';
import 'league_formula_section.dart';
import 'league_grants_section.dart';
import 'league_moves_section.dart';
import 'league_providers.dart';
import 'league_refresh.dart';
import 'league_table_section.dart';
import 'league_tasks_section.dart';
import 'league_transition.dart';
import 'league_view.dart';
import 'standing_card.dart';
import 'transition_status.dart';

class LeagueStudentBody extends ConsumerWidget {
  const LeagueStudentBody({required this.view, super.key});

  final LeagueView view;

  Future<void> _display(BuildContext context, WidgetRef ref) async {
    final standing = view.me!;
    final chosen = await LeagueDisplayDialog.show(context, standing);
    if (chosen == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(leagueApiProvider).display(view.id, chosen.$1, chosen.$2),
    );
    ref.refreshLeague(view.kind, view.originId);
  }

  Future<void> _dispute(
    BuildContext context,
    WidgetRef ref, {
    LeagueEntry? entry,
    LeagueTransition? move,
  }) async {
    final l10n = context.l10n;
    final text = await TeamTextDialog.show(
      context,
      title: l10n.leagueDispute,
      label: l10n.leagueDisputeText,
      hint: entry?.title,
      action: l10n.teamLeaveSend,
      lines: 3,
    );
    if (text == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref
          .read(leagueApiProvider)
          .dispute(view.id, text, entryId: entry?.id, transitionId: move?.id),
    );
    ref.refreshLeague(view.kind, view.originId);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final me = view.me;
    return SpacedColumn(
      children: [
        if (me != null) ...[
          StandingCard(
            view: view,
            standing: me,
            onDisplay: () => _display(context, ref),
          ),
          LeagueTableSection(
            division: view.division(me.division),
            rows: view.table,
          ),
        ],
        LeagueTasksSection(tasks: view.tasks),
        LeagueEntriesSection(
          entries: view.entries,
          onDispute: (entry) => _dispute(context, ref, entry: entry),
        ),
        if (view.transitions.isNotEmpty)
          LeagueMovesSection(
            view: view,
            moves: view.transitions,
            trailing: (move) => move.status == TransitionStatus.pending
                ? const SizedBox.shrink()
                : IconButton(
                    tooltip: l10n.leagueDispute,
                    icon: const Icon(Icons.report_outlined),
                    onPressed: () => _dispute(context, ref, move: move),
                  ),
          ),
        if (view.disputes.isNotEmpty)
          LeagueDisputesSection(disputes: view.disputes, entries: view.entries),
        if (view.grants.isNotEmpty) LeagueGrantsSection(grants: view.grants),
        LeagueFormulaSection(formula: view.formula),
      ],
    );
  }
}
