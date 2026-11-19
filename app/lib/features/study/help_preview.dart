import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'help_summary_view.dart';
import 'study_providers.dart';

class HelpPreview extends ConsumerWidget {
  const HelpPreview({required this.assignmentId, super.key});

  final String assignmentId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final summary = ref.watch(helpSummaryProvider(assignmentId)).value;
    final muted = Theme.of(context).textTheme.bodySmall!
        .copyWith(color: context.palette.muted);
    return LodestarCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(context.l10n.helpSummaryAttached, style: muted),
          if (summary != null) ...[
            const SizedBox(height: Spacing.s),
            HelpSummaryView(summary: summary),
          ],
        ],
      ),
    );
  }
}
