import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/components/spaced_column.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import '../team/team_text_dialog.dart';
import 'adjust_dialog.dart';
import 'award_picker.dart';
import 'league_entries_section.dart';
import 'league_grants_section.dart';
import 'league_moves_section.dart';
import 'league_providers.dart';
import 'league_refresh.dart';
import 'league_view.dart';
import 'player_details.dart';
import 'standing_card.dart';

class LeaguePlayerScreen extends ConsumerWidget {
  const LeaguePlayerScreen({
    required this.kind,
    required this.originId,
    required this.userId,
    super.key,
  });

  final String kind;
  final String originId;
  final String userId;

  void _refresh(WidgetRef ref, LeagueView view) {
    ref
      ..invalidate(playerProvider((view.id, userId)))
      ..refreshLeague(kind, originId);
  }

  Future<void> _adjust(
    BuildContext context,
    WidgetRef ref,
    LeagueView view,
  ) async {
    final input = await AdjustDialog.show(context);
    if (input == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref
          .read(leagueApiProvider)
          .adjust(view.id, userId, input.$1, input.$2),
    );
    _refresh(ref, view);
  }

  Future<void> _grant(
    BuildContext context,
    WidgetRef ref,
    LeagueView view,
  ) async {
    final l10n = context.l10n;
    final awards = await ref.read(awardsProvider.future);
    if (!context.mounted) {
      return;
    }
    final award = await AwardPicker.show(context, [
      for (final award in awards)
        if (award.left) award,
    ]);
    if (award == null || !context.mounted) {
      return;
    }
    final reason = await TeamTextDialog.show(
      context,
      title: award.title,
      label: l10n.leagueGrantReason,
      action: l10n.leagueGrant,
      required: false,
    );
    if (reason == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () =>
          ref.read(leagueApiProvider).grant(view.id, award.id, userId, reason),
    );
    ref.invalidate(awardsProvider);
    _refresh(ref, view);
  }

  List<Widget> _body(
    BuildContext context,
    WidgetRef ref,
    LeagueView view,
    PlayerDetails player,
  ) {
    final l10n = context.l10n;
    return [
      Text(
        player.person.label,
        style: Theme.of(context).textTheme.headlineSmall,
      ),
      const SizedBox(height: Spacing.l),
      SpacedColumn(
        children: [
          StandingCard(view: view, standing: player.standing),
          Wrap(
            spacing: Spacing.s,
            children: [
              OutlinedButton.icon(
                key: const ValueKey('league-adjust'),
                onPressed: () => _adjust(context, ref, view),
                icon: const Icon(Icons.exposure),
                label: Text(l10n.leagueAdjust),
              ),
              OutlinedButton.icon(
                key: const ValueKey('league-grant'),
                onPressed: () => _grant(context, ref, view),
                icon: const Icon(Icons.emoji_events_outlined),
                label: Text(l10n.leagueGrant),
              ),
            ],
          ),
          LeagueEntriesSection(entries: player.entries),
          if (player.transitions.isNotEmpty)
            LeagueMovesSection(view: view, moves: player.transitions),
          if (player.grants.isNotEmpty)
            LeagueGrantsSection(grants: player.grants),
        ],
      ),
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final view = ref.watch(leagueProvider((kind, originId))).value;
    final player = view == null
        ? const AsyncValue<PlayerDetails>.loading()
        : ref.watch(playerProvider((view.id, userId)));
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.leagueTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: player,
          onRefresh: () async {
            if (view != null) {
              _refresh(ref, view);
            }
          },
          builder: (player) => _body(context, ref, view!, player),
        ),
      ),
    );
  }
}
