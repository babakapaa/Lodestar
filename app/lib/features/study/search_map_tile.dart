import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'assignment.dart';
import 'map_labels.dart';
import 'map_rules.dart';
import 'study_providers.dart';

class SearchMapTile extends ConsumerWidget {
  const SearchMapTile({required this.assignment, super.key});

  final Assignment assignment;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final labels = MapLabels(l10n);
    final rules = assignment.searchMap;
    final map = ref.watch(searchMapProvider(assignment.id)).value?.value.map;
    final gaps = map?.gaps(rules ?? MapRules.none) ?? const [];
    final palette = context.palette;
    return LodestarCard(
      padding: EdgeInsets.zero,
      child: ListTile(
        contentPadding: LodestarCard.tilePadding,
        leading: Icon(
          map != null && gaps.isEmpty
              ? Icons.check_circle_outline
              : Icons.travel_explore,
          color: map != null && gaps.isEmpty ? palette.accent : null,
        ),
        title: Text(l10n.mapOpen),
        subtitle: Text(
          [
            labels.rules(rules),
            if (map != null && gaps.isNotEmpty)
              labels.gap(gaps.first, map, rules ?? MapRules.none)
            else if (map != null)
              l10n.mapReady,
          ].join('\n'),
        ),
        trailing: const Icon(Icons.chevron_right),
        onTap: () => context.push(AppRoutes.searchMap(assignment.id)),
      ),
    );
  }
}
