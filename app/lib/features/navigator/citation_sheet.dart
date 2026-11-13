import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import 'citation.dart';

class CitationSheet extends StatelessWidget {
  const CitationSheet({required this.citation, super.key});

  final Citation citation;

  static Future<void> show(BuildContext context, Citation citation) {
    return showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      builder: (context) => CitationSheet(citation: citation),
    );
  }

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;
    final muted = textTheme.bodyMedium!.copyWith(color: context.palette.muted);
    return SafeArea(
      child: SingleChildScrollView(
        padding: const EdgeInsets.fromLTRB(Spacing.l, 0, Spacing.l, Spacing.xl),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '[${citation.n}] ${citation.title}',
              style: textTheme.titleMedium,
            ),
            if (citation.reference.isNotEmpty) ...[
              const SizedBox(height: Spacing.s),
              SelectableText(citation.reference, style: muted),
            ],
            if (citation.url.isNotEmpty) ...[
              const SizedBox(height: Spacing.s),
              SelectableText(citation.url, style: muted),
            ],
            const SizedBox(height: Spacing.l),
            SelectableText(citation.text),
          ],
        ),
      ),
    );
  }
}
