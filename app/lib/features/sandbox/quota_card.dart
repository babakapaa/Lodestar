import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'sandbox_labels_context.dart';
import 'sandbox_quota.dart';

class QuotaCard extends StatelessWidget {
  const QuotaCard({required this.quota, super.key});

  final SandboxQuota quota;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    final theme = Theme.of(context);
    final muted = theme.textTheme.bodyMedium!.copyWith(
      color: context.palette.muted,
    );
    return LodestarCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(l10n.sandboxBalance, style: muted),
          const SizedBox(height: Spacing.xs),
          Text(
            labels.hours(quota.available),
            key: const ValueKey('sandbox-balance'),
            style: theme.textTheme.headlineMedium,
          ),
          const SizedBox(height: Spacing.s),
          LinearProgressIndicator(value: quota.share),
          const SizedBox(height: Spacing.s),
          Text(
            l10n.sandboxBalanceOf(
              labels.hours(quota.hours),
              labels.hours(quota.granted),
            ),
            style: muted,
          ),
          if (quota.reserved > 0)
            Text(
              l10n.sandboxReserved(labels.hours(quota.reserved)),
              style: muted,
            ),
          const SizedBox(height: Spacing.s),
          Text(
            quota.forecast >= 0
                ? l10n.sandboxForecast(labels.hours(quota.forecast))
                : l10n.sandboxForecastShort,
            key: const ValueKey('sandbox-forecast'),
          ),
          ExpansionTile(
            tilePadding: EdgeInsets.zero,
            title: Text(l10n.sandboxWhereFrom),
            children: [
              for (final line in quota.lines)
                ListTile(
                  dense: true,
                  contentPadding: EdgeInsets.zero,
                  title: Text(labels.source(line.source, line.title)),
                  trailing: Text(
                    '${line.hours > 0 ? '+' : ''}${labels.hours(line.hours)}',
                  ),
                ),
            ],
          ),
        ],
      ),
    );
  }
}
