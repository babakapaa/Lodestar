import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'sandbox_providers.dart';
import 'sandbox_term.dart';
import 'term_dialog.dart';

class TermsSection extends ConsumerWidget {
  const TermsSection({super.key});

  Future<void> _edit(
    BuildContext context,
    WidgetRef ref, [
    SandboxTerm? term,
  ]) async {
    final value = await TermDialog.show(context, term);
    if (value == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(sandboxApiProvider).saveTerm(value),
    );
    ref.invalidate(sandboxTermsProvider);
  }

  Future<void> _delete(
    BuildContext context,
    WidgetRef ref,
    SandboxTerm term,
  ) async {
    await TeamAction.run(
      context,
      () => ref.read(sandboxApiProvider).deleteTerm(term.id),
    );
    ref.invalidate(sandboxTermsProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final terms = ref.watch(sandboxTermsProvider).value ?? const [];
    return Section(
      title: l10n.sandboxTerms,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          for (final term in terms)
            LodestarCard(
              key: ValueKey('term-${term.id}'),
              child: ListTile(
                contentPadding: EdgeInsets.zero,
                title: Text(term.name),
                subtitle: Text('${term.startsOn} – ${term.endsOn}'),
                onTap: () => _edit(context, ref, term),
                trailing: IconButton(
                  tooltip: l10n.sandboxTermDelete,
                  icon: const Icon(Icons.delete_outline),
                  onPressed: () => _delete(context, ref, term),
                ),
              ),
            ),
          TextButton.icon(
            key: const ValueKey('term-new'),
            onPressed: () => _edit(context, ref),
            icon: const Icon(Icons.add),
            label: Text(l10n.sandboxTermNew),
          ),
        ],
      ),
    );
  }
}
