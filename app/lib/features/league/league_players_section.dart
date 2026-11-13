import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'league_labels_context.dart';
import 'league_view.dart';

class LeaguePlayersSection extends StatelessWidget {
  const LeaguePlayersSection({required this.view, super.key});

  final LeagueView view;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.leagueLabels;
    return Section(
      title: l10n.leaguePlayers(view.players.length),
      child: Column(
        children: [
          for (final player in view.players)
            ListTile(
              key: ValueKey('player-${player.person.userId}'),
              contentPadding: EdgeInsets.zero,
              title: Text(player.person.label),
              subtitle: Text(
                [
                  view.division(player.division),
                  labels.growth(player.growth),
                  if (player.topGrowth) l10n.leagueTopGrowth,
                ].join(', '),
              ),
              trailing: Text('${player.points}'),
              onTap: () => context.push(
                AppRoutes.leaguePlayer(
                  view.kind,
                  view.originId,
                  player.person.userId,
                ),
              ),
            ),
        ],
      ),
    );
  }
}
