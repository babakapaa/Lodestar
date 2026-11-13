import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class SampleDialog extends StatefulWidget {
  const SampleDialog({super.key});

  static Future<int?> show(BuildContext context) => showDialog<int>(
    context: context,
    builder: (context) => const SampleDialog(),
  );

  @override
  State<SampleDialog> createState() => _SampleDialogState();
}

class _SampleDialogState extends State<SampleDialog> {
  final _count = TextEditingController(text: '3');

  @override
  void dispose() {
    _count.dispose();
    super.dispose();
  }

  void _done() {
    final count = int.tryParse(_count.text.trim());
    if (count != null && count > 0) {
      Navigator.of(context).pop(count);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.defenseSample),
      content: SizedBox(
        width: 380,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(l10n.defenseSampleHint),
            TextField(
              key: const ValueKey('sample-count'),
              controller: _count,
              autofocus: true,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(labelText: l10n.defenseSampleCount),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(onPressed: _done, child: Text(l10n.defenseSampleChoose)),
      ],
    );
  }
}
