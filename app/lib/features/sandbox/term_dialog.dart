import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'sandbox_term.dart';

class TermDialog extends StatefulWidget {
  const TermDialog({this.value, super.key});

  final SandboxTerm? value;

  static Future<SandboxTerm?> show(
    BuildContext context, [
    SandboxTerm? value,
  ]) => showDialog<SandboxTerm>(
    context: context,
    builder: (context) => TermDialog(value: value),
  );

  @override
  State<TermDialog> createState() => _TermDialogState();
}

class _TermDialogState extends State<TermDialog> {
  late final _name = TextEditingController(text: widget.value?.name ?? '');
  late final _from = TextEditingController(text: widget.value?.startsOn ?? '');
  late final _to = TextEditingController(text: widget.value?.endsOn ?? '');

  @override
  void dispose() {
    _name.dispose();
    _from.dispose();
    _to.dispose();
    super.dispose();
  }

  void _done() {
    final from = DateTime.tryParse(_from.text.trim());
    final to = DateTime.tryParse(_to.text.trim());
    if (_name.text.trim().isEmpty ||
        from == null ||
        to == null ||
        !to.isAfter(from)) {
      return;
    }
    Navigator.of(context).pop(
      SandboxTerm(
        id: widget.value?.id ?? '',
        name: _name.text.trim(),
        startsOn: _from.text.trim(),
        endsOn: _to.text.trim(),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(
        widget.value == null ? l10n.sandboxTermNew : l10n.sandboxTermEdit,
      ),
      content: SizedBox(
        width: 380,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              key: const ValueKey('term-name'),
              controller: _name,
              decoration: InputDecoration(labelText: l10n.sandboxTermName),
            ),
            TextField(
              key: const ValueKey('term-from'),
              controller: _from,
              decoration: InputDecoration(
                labelText: l10n.sandboxTermFrom,
                helperText: l10n.sandboxDateHint,
              ),
            ),
            TextField(
              key: const ValueKey('term-to'),
              controller: _to,
              decoration: InputDecoration(
                labelText: l10n.sandboxTermTo,
                helperText: l10n.sandboxDateHint,
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
        FilledButton(
          key: const ValueKey('term-save'),
          onPressed: _done,
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
