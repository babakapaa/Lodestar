import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'gpu_class.dart';
import 'gpu_node.dart';
import 'node_dialog.dart';
import 'node_token_dialog.dart';
import 'sandbox_providers.dart';

class NodesSection extends ConsumerWidget {
  const NodesSection({super.key});

  Future<void> _create(BuildContext context, WidgetRef ref) async {
    final classes = [
      for (final c in ref.read(gpuClassesProvider).value ?? const <GpuClass>[])
        if (!c.archived) c,
    ];
    final input = await NodeDialog.show(context, classes);
    if (input == null || !context.mounted) {
      return;
    }
    GpuNode? created;
    await TeamAction.run(context, () async {
      created = await ref
          .read(sandboxApiProvider)
          .createNode(
            name: input.name,
            classId: input.classId,
            gpus: input.gpus,
          );
    });
    ref.invalidate(gpuNodesProvider);
    final token = created?.token;
    if (token != null && context.mounted) {
      await NodeTokenDialog.show(context, token);
    }
  }

  Future<void> _delete(
    BuildContext context,
    WidgetRef ref,
    GpuNode node,
  ) async {
    await TeamAction.run(
      context,
      () => ref.read(sandboxApiProvider).deleteNode(node.id),
    );
    ref.invalidate(gpuNodesProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final nodes = ref.watch(gpuNodesProvider).value ?? const [];
    return Section(
      title: l10n.sandboxNodes,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(l10n.sandboxNodesHint),
          for (final n in nodes)
            LodestarCard(
              key: ValueKey('node-${n.id}'),
              child: ListTile(
                contentPadding: EdgeInsets.zero,
                leading: Icon(
                  n.online ? Icons.dns_outlined : Icons.cloud_off_outlined,
                ),
                title: Text(n.name),
                subtitle: Text(
                  [
                    n.className,
                    l10n.sandboxNodeBusy(n.busy, n.gpus),
                    if (n.online) l10n.sandboxOnline else l10n.sandboxOffline,
                  ].join(', '),
                ),
                trailing: IconButton(
                  tooltip: l10n.sandboxNodeDelete,
                  icon: const Icon(Icons.delete_outline),
                  onPressed: () => _delete(context, ref, n),
                ),
              ),
            ),
          TextButton.icon(
            key: const ValueKey('node-new'),
            onPressed: () => _create(context, ref),
            icon: const Icon(Icons.add),
            label: Text(l10n.sandboxNodeNew),
          ),
        ],
      ),
    );
  }
}
