import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/async_list_view.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'stats_period.dart';
import 'stats_providers.dart';
import 'stats_report.dart';
import 'unit_stats_card.dart';

class UnitStatsScreen extends ConsumerStatefulWidget {
  const UnitStatsScreen({this.parentId, this.title, super.key});

  final String? parentId;
  final String? title;

  @override
  ConsumerState<UnitStatsScreen> createState() => _UnitStatsScreenState();
}

class _UnitStatsScreenState extends ConsumerState<UnitStatsScreen> {
  int _days = 30;

  Future<StatsReport> _refresh() => ref.refresh(
    unitStatsProvider((parentId: widget.parentId, days: _days)).future,
  );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final report = ref.watch(
      unitStatsProvider((parentId: widget.parentId, days: _days)),
    );
    return Scaffold(
      appBar: AppBar(title: Text(widget.title ?? l10n.statsTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: report,
          onRefresh: _refresh,
          builder: (report) => [
            StatsPeriod(
              days: _days,
              onChanged: (days) => setState(() => _days = days),
            ),
            const SizedBox(height: Spacing.l),
            Text(
              report.full
                  ? l10n.statsFull
                  : l10n.statsPrivacy(report.threshold),
            ),
            const SizedBox(height: Spacing.l),
            if (report.units.isEmpty) Text(l10n.statsEmpty),
            for (final unit in report.units)
              UnitStatsCard(
                key: ValueKey('unit-stats-${unit.unitId}'),
                unit: unit,
                threshold: report.threshold,
                onTap: unit.hidden
                    ? null
                    : () => context.push(
                        AppRoutes.unitStats(unit.unitId, unit.name),
                      ),
              ),
          ],
        ),
      ),
    );
  }
}
