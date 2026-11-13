import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'sandbox_providers.dart';
import 'sandbox_run.dart';

class HeldRunsSection extends ConsumerWidget {
  const HeldRunsSection({super.key});

  Future<void> _decide(
    BuildContext context,
    WidgetRef ref,
    SandboxRun run, {
    required bool release,
  }) async {
    await TeamAction.run(
      context,
      () => ref.read(sandboxApiProvider).decideHold(run.id, release: release),
    );
    ref.invalidate(heldRunsProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final held = ref.watch(heldRunsProvider).value ?? const [];
    return Section(
      title: l10n.sandboxHeldTitle,
      child: held.isEmpty
          ? Text(l10n.sandboxHeldNone)
          : Column(
              children: [
                for (final run in held)
                  LodestarCard(
                    key: ValueKey('held-${run.id}'),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        ListTile(
                          contentPadding: EdgeInsets.zero,
                          leading: const Icon(Icons.pause_circle_outline),
                          title: Text('${run.displayName}, ${run.title}'),
                          subtitle: Text(l10n.sandboxHeldWhy(run.holdNote)),
                        ),
                        OverflowBar(
                          alignment: MainAxisAlignment.end,
                          children: [
                            TextButton(
                              key: ValueKey('held-stop-${run.id}'),
                              onPressed: () =>
                                  _decide(context, ref, run, release: false),
                              child: Text(l10n.sandboxStop),
                            ),
                            FilledButton(
                              key: ValueKey('held-release-${run.id}'),
                              onPressed: () =>
                                  _decide(context, ref, run, release: true),
                              child: Text(l10n.sandboxRelease),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
              ],
            ),
    );
  }
}
