import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../league/league_providers.dart';
import 'share_bar.dart';

class DivisionsSection extends ConsumerWidget {
  const DivisionsSection({required this.courseId, super.key});

  final String courseId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final view = ref.watch(leagueProvider(('study.course', courseId))).value;
    if (view == null || view.players.isEmpty || view.divisions.isEmpty) {
      return const SizedBox.shrink();
    }
    final l10n = context.l10n;
    final counts = List.filled(view.divisions.length, 0);
    for (final p in view.players) {
      if (p.division >= 0 && p.division < counts.length) {
        counts[p.division]++;
      }
    }
    return Padding(
      padding: const EdgeInsets.only(top: Spacing.xl),
      child: Section(
        title: l10n.insightsDivisions,
        child: Column(
          children: [
            for (var i = counts.length - 1; i >= 0; i--)
              ShareBar(
                key: ValueKey('division-$i'),
                label: view.division(i),
                value: l10n.insightsStudents(counts[i]),
                share: counts[i] / view.players.length,
              ),
          ],
        ),
      ),
    );
  }
}
