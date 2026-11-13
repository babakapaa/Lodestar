import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class EditMessageDialog extends StatefulWidget {
  const EditMessageDialog({required this.text, super.key});

  final String text;

  static Future<String?> show(BuildContext context, String text) =>
      showDialog<String>(
        context: context,
        builder: (context) => EditMessageDialog(text: text),
      );

  @override
  State<EditMessageDialog> createState() => _EditMessageDialogState();
}

class _EditMessageDialogState extends State<EditMessageDialog> {
  late final _text = TextEditingController(text: widget.text);

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.chatsEditTitle),
      content: SizedBox(
        width: 480,
        child: TextField(
          key: const ValueKey('chat-edit'),
          controller: _text,
          autofocus: true,
          minLines: 1,
          maxLines: 8,
          maxLength: 4000,
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(
          onPressed: () {
            final text = _text.text.trim();
            if (text.isNotEmpty) {
              Navigator.of(context).pop(text);
            }
          },
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
