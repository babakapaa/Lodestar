import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'gpu_class.dart';
import 'node_input.dart';

class NodeDialog extends StatefulWidget {
  const NodeDialog({required this.classes, super.key});

  final List<GpuClass> classes;

  static Future<NodeInput?> show(
    BuildContext context,
    List<GpuClass> classes,
  ) => showDialog<NodeInput>(
    context: context,
    builder: (context) => NodeDialog(classes: classes),
  );

  @override
  State<NodeDialog> createState() => _NodeDialogState();
}

class _NodeDialogState extends State<NodeDialog> {
  final _name = TextEditingController();
  final _gpus = TextEditingController(text: '1');
  late String? _classId = widget.classes.firstOrNull?.id;

  @override
  void dispose() {
    _name.dispose();
    _gpus.dispose();
    super.dispose();
  }

  void _done() {
    final gpus = int.tryParse(_gpus.text.trim());
    final classId = _classId;
    if (_name.text.trim().isEmpty || gpus == null || classId == null) {
      return;
    }
    Navigator.of(context)
        .pop(NodeInput(name: _name.text.trim(), classId: classId, gpus: gpus));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.sandboxNodeNew),
      content: SizedBox(
        width: 380,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              key: const ValueKey('node-name'),
              controller: _name,
              decoration: InputDecoration(labelText: l10n.sandboxNodeName),
            ),
            DropdownButtonFormField<String>(
              key: const ValueKey('node-class'),
              initialValue: _classId,
              isExpanded: true,
              decoration: InputDecoration(labelText: l10n.sandboxFactClass),
              items: [
                for (final c in widget.classes)
                  DropdownMenuItem(value: c.id, child: Text(c.name)),
              ],
              onChanged: (id) => setState(() => _classId = id),
            ),
            TextField(
              key: const ValueKey('node-gpus'),
              controller: _gpus,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(labelText: l10n.sandboxNodeGpus),
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
          key: const ValueKey('node-save'),
          onPressed: _done,
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
