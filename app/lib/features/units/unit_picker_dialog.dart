import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'unit.dart';
import 'unit_kinds.dart';
import 'unit_target.dart';
import 'units_providers.dart';
import 'units_query.dart';

class UnitPickerDialog extends ConsumerStatefulWidget {
  const UnitPickerDialog({required this.moving, super.key});

  final Unit moving;

  static Future<UnitTarget?> show(BuildContext context, Unit moving) =>
      showDialog<UnitTarget>(
        context: context,
        builder: (context) => UnitPickerDialog(moving: moving),
      );

  @override
  ConsumerState<UnitPickerDialog> createState() => _UnitPickerDialogState();
}

class _UnitPickerDialogState extends ConsumerState<UnitPickerDialog> {
  final List<Unit> _trail = [];

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final here = _trail.isEmpty ? null : _trail.last;
    final children = ref.watch(
      unitChildrenProvider(UnitsQuery(parentId: here?.id, archived: false)),
    );
    return AlertDialog(
      title: Text(here?.name ?? l10n.unitsRoot),
      content: SizedBox(
        width: 420,
        height: 360,
        child: switch (children) {
          AsyncData(:final value) => ListView(
            children: [
              if (here != null)
                ListTile(
                  leading: const Icon(Icons.arrow_upward),
                  title: Text(l10n.unitsUp),
                  onTap: () => setState(_trail.removeLast),
                ),
              for (final unit in value)
                if (unit.id != widget.moving.id)
                  ListTile(
                    title: Text(unit.name),
                    subtitle: Text(UnitKinds.label(l10n, unit.kind)),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => setState(() => _trail.add(unit)),
                  ),
            ],
          ),
          AsyncError() => Text(l10n.errorGeneric),
          _ => const Center(child: CircularProgressIndicator()),
        },
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        TextButton(
          onPressed: here?.id == widget.moving.parentId
              ? null
              : () => Navigator.of(context).pop(UnitTarget(here?.id)),
          child: Text(l10n.unitsMoveHere),
        ),
      ],
    );
  }
}
