import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'league_labels_context.dart';
import 'league_view.dart';
import 'standing.dart';

class StandingCard extends StatelessWidget {
  const StandingCard({
    required this.view,
    required this.standing,
    this.onDisplay,
    super.key,
  });

  final LeagueView view;
  final Standing standing;
  final VoidCallback? onDisplay;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.leagueLabels;
    final textTheme = Theme.of(context).textTheme;
    final muted = textTheme.bodyMedium!.copyWith(color: context.palette.muted);
    final display = onDisplay;
    return LodestarCard(
      key: const ValueKey('standing'),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(view.division(standing.division), style: textTheme.titleLarge),
          const SizedBox(height: Spacing.s),
          Text(
            l10n.leaguePoints(standing.points),
            style: textTheme.titleMedium,
          ),
          Text(
            labels.toNext(standing, view.division(standing.division + 1)),
            style: muted,
          ),
          Text(labels.growth(standing.growth), style: muted),
          if (standing.gpuHours > 0)
            Text(l10n.leagueGpuHours(standing.gpuHours), style: muted),
          if (display != null) ...[
            const SizedBox(height: Spacing.s),
            TextButton.icon(
              key: const ValueKey('league-display'),
              onPressed: display,
              icon: const Icon(Icons.visibility_outlined),
              label: Text(l10n.leagueShownAs(labels.display(standing.display))),
            ),
          ],
        ],
      ),
    );
  }
}
