import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'map_source.dart';
import 'search_map.dart';
import 'search_map_diagram.dart';

class SearchMapView extends StatelessWidget {
  const SearchMapView({required this.map, super.key});

  final SearchMap map;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final textTheme = Theme.of(context).textTheme;
    final muted = textTheme.bodySmall!.copyWith(color: context.palette.muted);
    List<Widget> sources(String title, List<MapSource> list, String why) => [
      const SizedBox(height: Spacing.m),
      Text(title, style: textTheme.titleSmall),
      if (list.isEmpty) Text(l10n.mapNoneHere, style: muted),
      for (final (index, source) in list.indexed) ...[
        const SizedBox(height: Spacing.s),
        SelectableText('${index + 1}. ${source.reference}'),
        for (final (label, value) in [
          (l10n.mapFound, source.found),
          (why, source.why),
          (l10n.mapTaken, source.taken),
        ])
          if (value.trim().isNotEmpty) Text('$label: $value', style: muted),
      ],
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(l10n.mapTitle, style: textTheme.titleMedium),
        ...sources(l10n.mapUsed, map.used, l10n.mapWhyUsed),
        ...sources(l10n.mapRejected, map.rejected, l10n.mapWhyRejected),
        const SizedBox(height: Spacing.m),
        Text(l10n.mapPath, style: textTheme.titleSmall),
        SelectableText(map.path),
        for (final query in map.queries) Text('• $query', style: muted),
        SearchMapDiagram(map: map),
      ],
    );
  }
}
