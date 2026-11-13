import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'search_map.dart';
import 'search_map_painter.dart';

class SearchMapDiagram extends StatelessWidget {
  const SearchMapDiagram({required this.map, super.key});

  final SearchMap map;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return Semantics(
      label: l10n.mapSchemeLabel(map.used.length, map.rejected.length),
      excludeSemantics: true,
      child: SizedBox(
        height: 320,
        child: CustomPaint(
          size: Size.infinite,
          painter: SearchMapPainter(
            map: map,
            center: l10n.mapWork,
            palette: context.palette,
            labelStyle: Theme.of(context).textTheme.labelSmall!,
          ),
        ),
      ),
    );
  }
}
