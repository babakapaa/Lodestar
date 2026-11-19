import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'help_summary.dart';

class HelpSummaryView extends StatelessWidget {
  const HelpSummaryView({required this.summary, super.key});

  final HelpSummary summary;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final textTheme = Theme.of(context).textTheme;
    final muted = textTheme.bodyMedium!.copyWith(color: context.palette.muted);
    final steps = summary.steps;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(l10n.helpSummaryTitle, style: textTheme.labelLarge),
        const SizedBox(height: Spacing.xs),
        if (summary.replies == 0)
          Text(l10n.helpSummaryNone, style: muted)
        else ...[
          Text(l10n.helpSummaryReplies(summary.replies)),
          Text(
            l10n.helpSummarySteps(
              steps.question,
              steps.direction,
              steps.source,
              steps.fragment,
              steps.example,
            ),
            style: muted,
          ),
          Text(l10n.helpSummaryRule, style: muted),
        ],
        if (summary.sources.isNotEmpty) ...[
          const SizedBox(height: Spacing.s),
          Text(l10n.helpSummarySources, style: textTheme.labelMedium),
          for (final source in summary.sources) Text(source, style: muted),
        ],
        if (summary.articles.isNotEmpty) ...[
          const SizedBox(height: Spacing.s),
          Text(l10n.helpSummaryArticles, style: textTheme.labelMedium),
          for (final article in summary.articles) Text(article, style: muted),
        ],
      ],
    );
  }
}
