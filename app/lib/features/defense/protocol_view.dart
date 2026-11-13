import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'defense.dart';
import 'defense_labels.dart';

class ProtocolView extends StatelessWidget {
  const ProtocolView({required this.defense, super.key});

  final Defense defense;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = DefenseLabels(l10n);
    final textTheme = Theme.of(context).textTheme;
    final muted = textTheme.bodySmall!.copyWith(color: context.palette.muted);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (defense.summary.isNotEmpty) ...[
          SelectableText(defense.summary),
          const SizedBox(height: Spacing.m),
        ],
        for (final (i, q) in defense.questions.indexed) ...[
          Text('${i + 1}. ${q.text}', style: textTheme.titleSmall),
          Text('${labels.kind(q.kind)}, ${labels.mark(q.mark)}', style: muted),
          if (q.note.isNotEmpty) Text(q.note),
          const SizedBox(height: Spacing.s),
        ],
      ],
    );
  }
}
