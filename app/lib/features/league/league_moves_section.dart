import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'league_labels_context.dart';
import 'league_transition.dart';
import 'league_view.dart';

class LeagueMovesSection extends StatelessWidget {
  const LeagueMovesSection({
    required this.view,
    required this.moves,
    this.trailing,
    this.title,
    super.key,
  });

  final LeagueView view;
  final List<LeagueTransition> moves;
  final Widget Function(LeagueTransition move)? trailing;
  final String? title;

  @override
  Widget build(BuildContext context) {
    final labels = context.leagueLabels;
    final action = trailing;
    return Section(
      title: title ?? context.l10n.leagueMoves,
      child: Column(
        children: [
          for (final move in moves)
            ListTile(
              key: ValueKey('move-${move.id}'),
              contentPadding: EdgeInsets.zero,
              leading: Icon(move.up ? Icons.trending_up : Icons.trending_down),
              title: Text(
                [
                  if (move.person case final person?) person.label,
                  labels.transition(
                    move,
                    view.division(move.from),
                    view.division(move.to),
                  ),
                ].join(': '),
              ),
              subtitle: Text(
                [
                  labels.transitionStatus(move.status),
                  labels.time(move.createdAt),
                ].join(', '),
              ),
              trailing: action?.call(move),
            ),
        ],
      ),
    );
  }
}
