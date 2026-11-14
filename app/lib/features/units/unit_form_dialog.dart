import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'unit.dart';
import 'unit_input.dart';
import 'unit_kinds.dart';

class UnitFormDialog extends StatefulWidget {
  const UnitFormDialog({this.unit, this.parentId, super.key});

  final Unit? unit;
  final String? parentId;

  static Future<UnitInput?> show(
    BuildContext context, {
    Unit? unit,
    String? parentId,
  }) => showDialog<UnitInput>(
    context: context,
    builder: (context) => UnitFormDialog(unit: unit, parentId: parentId),
  );

  @override
  State<UnitFormDialog> createState() => _UnitFormDialogState();
}

class _UnitFormDialogState extends State<UnitFormDialog> {
  late final _name = TextEditingController(text: widget.unit?.name ?? '');
  late final _code = TextEditingController(text: widget.unit?.code ?? '');
  late String _kind = widget.unit?.kind ?? UnitKinds.presets.last;

  @override
  void dispose() {
    _name.dispose();
    _code.dispose();
    super.dispose();
  }

  void _submit() {
    final name = _name.text.trim();
    if (name.isEmpty) {
      return;
    }
    Navigator.of(context).pop(
      UnitInput(
        parentId: widget.unit == null ? widget.parentId : null,
        kind: _kind,
        name: name,
        code: _code.text.trim(),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final kinds = {...UnitKinds.presets, _kind};
    return AlertDialog(
      title: Text(widget.unit == null ? l10n.unitsCreate : l10n.unitsEdit),
      content: SizedBox(
        width: 420,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            DropdownButtonFormField<String>(
              isExpanded: true,
              initialValue: _kind,
              decoration: InputDecoration(labelText: l10n.unitsKind),
              items: [
                for (final kind in kinds)
                  DropdownMenuItem(
                    value: kind,
                    child: Text(UnitKinds.label(l10n, kind)),
                  ),
              ],
              onChanged: (kind) => setState(() => _kind = kind ?? _kind),
            ),
            TextField(
              controller: _name,
              autofocus: true,
              maxLength: 200,
              decoration: InputDecoration(labelText: l10n.unitsName),
            ),
            TextField(
              controller: _code,
              maxLength: 100,
              decoration: InputDecoration(
                labelText: l10n.unitsCode,
                helperText: l10n.unitsCodeHint,
              ),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        TextButton(onPressed: _submit, child: Text(l10n.actionSave)),
      ],
    );
  }
}
