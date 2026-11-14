import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'sandbox_providers.dart';
import 'sandbox_template.dart';
import 'template_dialog.dart';

class TemplatesSection extends ConsumerWidget {
  const TemplatesSection({super.key});

  Future<void> _edit(
    BuildContext context,
    WidgetRef ref, [
    SandboxTemplate? value,
  ]) async {
    final next = await TemplateDialog.show(context, value);
    if (next == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(sandboxApiProvider).saveTemplate(next),
    );
    ref.invalidate(sandboxTemplatesProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final templates = ref.watch(sandboxTemplatesProvider).value ?? const [];
    return Section(
      title: l10n.sandboxTemplates,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          for (final t in templates)
            LodestarCard(
              key: ValueKey('template-${t.id}'),
              padding: EdgeInsets.zero,
              child: ListTile(
                contentPadding: LodestarCard.tilePadding,
                leading: Icon(
                  t.starter
                      ? Icons.school_outlined
                      : Icons.inventory_2_outlined,
                ),
                title: Text(
                  t.starter ? l10n.sandboxStarterName(t.name) : t.name,
                ),
                subtitle: Text(
                  [t.image, if (t.archived) l10n.sandboxArchived].join(', '),
                ),
                onTap: () => _edit(context, ref, t),
              ),
            ),
          TextButton.icon(
            key: const ValueKey('template-new'),
            onPressed: () => _edit(context, ref),
            icon: const Icon(Icons.add),
            label: Text(l10n.sandboxTemplateNew),
          ),
        ],
      ),
    );
  }
}
