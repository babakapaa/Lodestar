import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class ReturnNoteDialog extends StatefulWidget {
  const ReturnNoteDialog({super.key});

  static Future<String?> show(BuildContext context) => showDialog<String>(
    context: context,
    builder: (context) => const ReturnNoteDialog(),
  );

  @override
  State<ReturnNoteDialog> createState() => _ReturnNoteDialogState();
}

class _ReturnNoteDialogState extends State<ReturnNoteDialog> {
  final _note = TextEditingController();

  @override
  void dispose() {
    _note.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.setupReturnTitle),
      content: TextField(controller: _note, autofocus: true, maxLines: 3),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        TextButton(
          onPressed: () => Navigator.of(context).pop(_note.text.trim()),
          child: Text(l10n.setupReturn),
        ),
      ],
    );
  }
}
