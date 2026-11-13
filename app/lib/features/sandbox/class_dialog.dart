import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'gpu_class.dart';

class ClassDialog extends StatefulWidget {
  const ClassDialog({this.value, super.key});

  final GpuClass? value;

  static Future<GpuClass?> show(BuildContext context, [GpuClass? value]) =>
      showDialog<GpuClass>(
        context: context,
        builder: (context) => ClassDialog(value: value),
      );

  @override
  State<ClassDialog> createState() => _ClassDialogState();
}

class _ClassDialogState extends State<ClassDialog> {
  late final _name = TextEditingController(text: widget.value?.name ?? '');
  late final _memory = TextEditingController(
    text: widget.value == null ? '' : '${widget.value!.memoryGb}',
  );
  late final _rate = TextEditingController(
    text: widget.value == null ? '' : '${widget.value!.rate}',
  );
  late bool _archived = widget.value?.archived ?? false;

  @override
  void dispose() {
    _name.dispose();
    _memory.dispose();
    _rate.dispose();
    super.dispose();
  }

  void _done() {
    final memory = int.tryParse(_memory.text.trim());
    final rate = double.tryParse(_rate.text.trim().replaceAll(',', '.'));
    if (_name.text.trim().isEmpty || memory == null || rate == null) {
      return;
    }
    Navigator.of(context).pop(
      GpuClass(
        id: widget.value?.id ?? '',
        name: _name.text.trim(),
        memoryGb: memory,
        rate: rate,
        archived: _archived,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(
        widget.value == null ? l10n.sandboxClassNew : l10n.sandboxClassEdit,
      ),
      content: SizedBox(
        width: 380,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              key: const ValueKey('class-name'),
              controller: _name,
              decoration: InputDecoration(labelText: l10n.sandboxClassField),
            ),
            TextField(
              key: const ValueKey('class-memory'),
              controller: _memory,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(labelText: l10n.sandboxMemory),
            ),
            TextField(
              key: const ValueKey('class-rate'),
              controller: _rate,
              keyboardType: const TextInputType.numberWithOptions(
                decimal: true,
              ),
              decoration: InputDecoration(
                labelText: l10n.sandboxRate,
                helperText: l10n.sandboxRateHint,
              ),
            ),
            SwitchListTile(
              contentPadding: EdgeInsets.zero,
              title: Text(l10n.sandboxArchived),
              value: _archived,
              onChanged: (v) => setState(() => _archived = v),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(
          key: const ValueKey('class-save'),
          onPressed: _done,
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
