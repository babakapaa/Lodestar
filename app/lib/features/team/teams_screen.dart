import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/navigator_tile.dart';
import 'team_action.dart';
import 'team_card.dart';
import 'team_card_tile.dart';
import 'team_providers.dart';
import 'team_text_dialog.dart';

class TeamsScreen extends ConsumerWidget {
  const TeamsScreen({super.key});

  Future<void> _accept(
    BuildContext context,
    WidgetRef ref,
    TeamCard card,
  ) async {
    await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).accept(card.team.id),
    );
    ref
      ..invalidate(teamInvitationsProvider)
      ..invalidate(myTeamsProvider);
  }

  Future<void> _join(BuildContext context) async {
    final l10n = context.l10n;
    final code = await TeamTextDialog.show(
      context,
      title: l10n.teamJoinByCode,
      label: l10n.teamCode,
      action: l10n.teamOpen,
    );
    if (code != null && context.mounted) {
      await context.push(AppRoutes.teamJoin(code));
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final teams = ref.watch(myTeamsProvider);
    final invitations = ref.watch(teamInvitationsProvider).value ?? const [];
    final list = teams.value ?? const <TeamCard>[];
    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.teamsTitle),
        actions: [
          IconButton(
            tooltip: l10n.teamJoinByCode,
            icon: const Icon(Icons.qr_code_2_outlined),
            onPressed: () => _join(context),
          ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: RefreshIndicator(
          onRefresh: () {
            ref.invalidate(teamInvitationsProvider);
            return ref.refresh(myTeamsProvider.future);
          },
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              NavigatorTile(
                icon: Icons.badge_outlined,
                title: l10n.teamProfileTitle,
                subtitle: l10n.teamProfileHint,
                route: AppRoutes.teamProfile,
              ),
              if (invitations.isNotEmpty) ...[
                const SizedBox(height: Spacing.xl),
                Section(
                  title: l10n.teamInvitations,
                  child: Column(
                    children: [
                      for (final card in invitations)
                        TeamCardTile(
                          key: ValueKey('invitation-${card.team.id}'),
                          card: card,
                          trailing: TextButton(
                            onPressed: () => _accept(context, ref, card),
                            child: Text(l10n.teamAccept),
                          ),
                        ),
                    ],
                  ),
                ),
              ],
              const SizedBox(height: Spacing.xl),
              if (teams.isLoading && list.isEmpty)
                const Center(child: CircularProgressIndicator())
              else if (list.isEmpty)
                EmptyState(
                  icon: Icons.groups_outlined,
                  title: l10n.teamsEmpty,
                  body: l10n.teamsEmptyBody,
                )
              else
                Section(
                  title: l10n.teamsMine,
                  child: Column(
                    children: [
                      for (final card in list)
                        TeamCardTile(key: ValueKey(card.team.id), card: card),
                    ],
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
