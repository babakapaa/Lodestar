import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'navigator_providers.dart';
import 'support_card.dart';

class SupportSheet extends ConsumerWidget {
  const SupportSheet({super.key});

  static Future<void> show(BuildContext context) {
    return showModalBottomSheet<void>(
      context: context,
      showDragHandle: true,
      builder: (context) => const SupportSheet(),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final contacts = ref.watch(supportProvider);
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(Spacing.l, 0, Spacing.l, Spacing.xl),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              l10n.supportTitle,
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: Spacing.s),
            Text(l10n.supportBody),
            const SizedBox(height: Spacing.l),
            switch (contacts) {
              AsyncData(:final value) when value.isNotEmpty => SupportCard(
                contacts: value,
              ),
              AsyncLoading() => const Center(
                child: CircularProgressIndicator(),
              ),
              _ => SupportCard(contacts: l10n.supportFallback),
            },
          ],
        ),
      ),
    );
  }
}
