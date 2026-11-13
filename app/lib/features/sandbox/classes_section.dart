import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'class_dialog.dart';
import 'gpu_class.dart';
import 'sandbox_labels_context.dart';
import 'sandbox_providers.dart';

class ClassesSection extends ConsumerWidget {
  const ClassesSection({super.key});

  Future<void> _edit(
    BuildContext context,
    WidgetRef ref, [
    GpuClass? value,
  ]) async {
    final next = await ClassDialog.show(context, value);
    if (next == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(sandboxApiProvider).saveClass(next),
    );
    ref.invalidate(gpuClassesProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    final classes = ref.watch(gpuClassesProvider).value ?? const [];
    return Section(
      title: l10n.sandboxClasses,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          for (final c in classes)
            LodestarCard(
              key: ValueKey('gpu-class-${c.id}'),
              child: ListTile(
                contentPadding: EdgeInsets.zero,
                leading: const Icon(Icons.memory),
                title: Text(l10n.sandboxClassName(c.name, c.memoryGb)),
                subtitle: Text(
                  [
                    l10n.sandboxRatePerHour(labels.hours(c.rate)),
                    if (c.archived) l10n.sandboxArchived,
                  ].join(', '),
                ),
                onTap: () => _edit(context, ref, c),
              ),
            ),
          TextButton.icon(
            key: const ValueKey('class-new'),
            onPressed: () => _edit(context, ref),
            icon: const Icon(Icons.add),
            label: Text(l10n.sandboxClassNew),
          ),
        ],
      ),
    );
  }
}
