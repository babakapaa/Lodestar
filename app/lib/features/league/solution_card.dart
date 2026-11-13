import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'league_solution.dart';
import 'task_details.dart';

class SolutionCard extends StatelessWidget {
  const SolutionCard({
    required this.details,
    required this.solution,
    this.onReview,
    this.onExplain,
    super.key,
  });

  final TaskDetails details;
  final LeagueSolution solution;
  final VoidCallback? onReview;
  final VoidCallback? onExplain;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final textTheme = Theme.of(context).textTheme;
    final muted = textTheme.bodyMedium!.copyWith(color: context.palette.muted);
    final similar = solution.similar;
    final review = onReview;
    final explain = onExplain;
    return LodestarCard(
      key: ValueKey('solution-${solution.id}'),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(solution.teamName, style: textTheme.titleMedium),
          Text(
            [for (final id in solution.members) details.nameOf(id)].join(', '),
            style: muted,
          ),
          const SizedBox(height: Spacing.s),
          if (solution.work.text.isNotEmpty) SelectableText(solution.work.text),
          for (final link in solution.work.links) SelectableText(link),
          if (similar != null) ...[
            const SizedBox(height: Spacing.m),
            NoticeBanner(
              icon: Icons.content_copy_outlined,
              message: l10n.leagueSimilar(
                similar.teamName,
                similar.taskTitle,
                ((solution.similarity ?? 0) * 100).round(),
              ),
            ),
            Text(similar.text, maxLines: 4, overflow: TextOverflow.ellipsis),
          ],
          if (solution.explanation.isNotEmpty) ...[
            const SizedBox(height: Spacing.s),
            Text(l10n.leagueExplanation(solution.explanation), style: muted),
          ],
          const SizedBox(height: Spacing.s),
          Text(
            solution.reviewed
                ? l10n.leagueReviewed(
                    solution.points ?? 0,
                    [for (final id in solution.passed) details.nameOf(id)]
                        .join(', '),
                  )
                : l10n.leagueNotReviewed,
            style: muted,
          ),
          if (solution.comment.isNotEmpty) Text(solution.comment),
          if (review != null || explain != null)
            Align(
              alignment: Alignment.centerLeft,
              child: review != null
                  ? TextButton(
                      key: ValueKey('review-${solution.id}'),
                      onPressed: review,
                      child: Text(l10n.leagueReview),
                    )
                  : TextButton(
                      key: ValueKey('explain-${solution.id}'),
                      onPressed: explain,
                      child: Text(l10n.leagueExplain),
                    ),
            ),
        ],
      ),
    );
  }
}
