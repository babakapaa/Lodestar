import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../insights/share_bar.dart';
import 'metric_labels.dart';
import 'unit_stats.dart';

class UnitStatsCard extends StatelessWidget {
  const UnitStatsCard({
    required this.unit,
    required this.threshold,
    this.onTap,
    super.key,
  });

  final UnitStats unit;
  final int threshold;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = MetricLabels(l10n);
    final people = unit.people;
    return Card(
      margin: const EdgeInsets.only(bottom: Spacing.m),
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(Spacing.l),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    child: Text(
                      unit.name,
                      style: Theme.of(context).textTheme.titleMedium,
                    ),
                  ),
                  if (people != null) Text(l10n.statsPeople(people)),
                  if (onTap != null) const Icon(Icons.chevron_right),
                ],
              ),
              const SizedBox(height: Spacing.s),
              if (unit.hidden)
                Text(l10n.statsHidden(threshold))
              else if (unit.metrics.isEmpty)
                Text(l10n.statsNoData)
              else
                for (final metric in unit.metrics)
                  if (labels.isShare(metric.key))
                    ShareBar(
                      label: labels.name(metric.key),
                      value: labels.value(metric),
                      share: metric.value,
                    )
                  else
                    Padding(
                      padding: const EdgeInsets.symmetric(vertical: Spacing.xs),
                      child: Row(
                        children: [
                          Expanded(child: Text(labels.name(metric.key))),
                          Text(labels.value(metric)),
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
