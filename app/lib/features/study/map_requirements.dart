import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'map_labels.dart';
import 'map_rules.dart';
import 'search_map.dart';

class MapRequirements extends StatelessWidget {
  const MapRequirements({required this.map, required this.rules, super.key});

  final SearchMap map;
  final MapRules? rules;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = MapLabels(l10n);
    final palette = context.palette;
    final applied = rules ?? MapRules.none;
    final gaps = map.gaps(applied);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(labels.rules(rules)),
        const SizedBox(height: Spacing.s),
        if (gaps.isEmpty)
          ListTile(
            contentPadding: EdgeInsets.zero,
            leading: Icon(Icons.check_circle_outline, color: palette.accent),
            title: Text(l10n.mapReady),
          )
        else
          for (final gap in gaps)
            ListTile(
              contentPadding: EdgeInsets.zero,
              leading: Icon(Icons.error_outline, color: palette.warn),
              title: Text(labels.gap(gap, map, applied)),
            ),
      ],
    );
  }
}
