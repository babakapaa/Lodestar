import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'division_draft.dart';

class DivisionFields extends StatelessWidget {
  const DivisionFields({
    required this.index,
    required this.draft,
    this.onRemove,
    super.key,
  });

  final int index;
  final DivisionDraft draft;
  final VoidCallback? onRemove;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final remove = onRemove;
    return Padding(
      padding: const EdgeInsets.only(bottom: Spacing.m),
      child: Row(
        children: [
          Expanded(
            flex: 3,
            child: TextField(
              key: ValueKey('division-name-$index'),
              controller: draft.name,
              decoration: InputDecoration(labelText: l10n.leagueDivisionName),
            ),
          ),
          const SizedBox(width: Spacing.s),
          Expanded(
            flex: 2,
            child: TextField(
              key: ValueKey('division-from-$index'),
              controller: draft.from,
              enabled: index > 0,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(labelText: l10n.leagueDivisionFrom),
            ),
          ),
          const SizedBox(width: Spacing.s),
          Expanded(
            flex: 2,
            child: TextField(
              key: ValueKey('division-gpu-$index'),
              controller: draft.gpuHours,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(labelText: l10n.leagueDivisionGpu),
            ),
          ),
          IconButton(
            tooltip: l10n.leagueDivisionRemove,
            icon: const Icon(Icons.remove_circle_outline),
            onPressed: remove,
          ),
        ],
      ),
    );
  }
}
