import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'outbox_tile.dart';
import 'study_providers.dart';

class OutboxSection extends ConsumerWidget {
  const OutboxSection({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final drafts = ref.watch(outboxProvider);
    if (drafts.isEmpty) {
      return const SizedBox.shrink();
    }
    return Padding(
      padding: const EdgeInsets.only(bottom: Spacing.xl),
      child: Section(
        title: context.l10n.studyOutboxTitle,
        child: Column(
          children: [
            for (final draft in drafts)
              OutboxTile(key: ValueKey(draft.clientId), draft: draft),
          ],
        ),
      ),
    );
  }
}
