import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../l10n/l10n_context.dart';
import 'defense_labels.dart';
import 'defense_row.dart';
import 'defense_status.dart';
import 'slot_text.dart';

class DefenseRowTile extends StatelessWidget {
  const DefenseRowTile({
    required this.assignmentId,
    required this.row,
    required this.sampled,
    required this.onSelect,
    super.key,
  });

  final String assignmentId;
  final DefenseRow row;
  final bool sampled;
  final ValueChanged<bool> onSelect;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final locale = Localizations.localeOf(context).toLanguageTag();
    final slot = row.slot;
    final status = row.status == DefenseStatus.planned
        ? (slot == null ? l10n.defenseNoSlot : SlotText.of(slot, locale))
        : DefenseLabels(l10n).status(row.status);
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: sampled
          ? Checkbox(
              key: ValueKey('select-${row.student.userId}'),
              value: row.selected,
              onChanged: row.status.finished
                  ? null
                  : (value) => onSelect(value ?? false),
            )
          : Icon(row.status.finished ? Icons.task_alt : Icons.person_outline),
      title: Text(row.student.label),
      subtitle: Text(status),
      trailing: const Icon(Icons.chevron_right),
      onTap: () =>
          context.push(AppRoutes.defense(assignmentId, row.student.userId)),
    );
  }
}
