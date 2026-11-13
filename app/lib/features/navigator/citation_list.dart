import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'citation.dart';
import 'citation_sheet.dart';

class CitationList extends StatelessWidget {
  const CitationList({required this.citations, super.key});

  final List<Citation> citations;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          context.l10n.navigatorSources,
          style: Theme.of(context).textTheme.labelLarge,
        ),
        for (final citation in citations)
          ListTile(
            dense: true,
            contentPadding: EdgeInsets.zero,
            leading: Text('[${citation.n}]'),
            title: Text(citation.title),
            subtitle: citation.reference.isEmpty
                ? null
                : Text(
                    citation.reference,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
            onTap: () => CitationSheet.show(context, citation),
          ),
      ],
    );
  }
}
