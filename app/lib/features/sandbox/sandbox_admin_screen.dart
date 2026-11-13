import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/spaced_column.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'classes_section.dart';
import 'held_runs_section.dart';
import 'nodes_section.dart';
import 'sandbox_providers.dart';
import 'templates_section.dart';
import 'terms_section.dart';

class SandboxAdminScreen extends ConsumerWidget {
  const SandboxAdminScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.sandboxAdminTitle)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 720),
            child: RefreshIndicator(
              onRefresh: () async {
                ref
                  ..invalidate(heldRunsProvider)
                  ..invalidate(sandboxTermsProvider)
                  ..invalidate(gpuClassesProvider)
                  ..invalidate(sandboxTemplatesProvider)
                  ..invalidate(gpuNodesProvider);
              },
              child: ListView(
                padding: const EdgeInsets.all(Spacing.l),
                children: [
                  Text(l10n.sandboxAdminHint),
                  const SizedBox(height: Spacing.xl),
                  const SpacedColumn(
                    children: [
                      HeldRunsSection(),
                      TermsSection(),
                      ClassesSection(),
                      TemplatesSection(),
                      NodesSection(),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
