import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class ChatTitleDialog extends StatefulWidget {
  const ChatTitleDialog({required this.heading, super.key});

  final String heading;

  static Future<String?> show(BuildContext context, String heading) =>
      showDialog<String>(
        context: context,
        builder: (context) => ChatTitleDialog(heading: heading),
      );

  @override
  State<ChatTitleDialog> createState() => _ChatTitleDialogState();
}

class _ChatTitleDialogState extends State<ChatTitleDialog> {
  final _title = TextEditingController();

  @override
  void dispose() {
    _title.dispose();
    super.dispose();
  }

  void _done() {
    final title = _title.text.trim();
    if (title.isNotEmpty) {
      Navigator.of(context).pop(title);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(widget.heading),
      content: TextField(
        key: const ValueKey('chat-title'),
        controller: _title,
        autofocus: true,
        maxLength: 200,
        decoration: InputDecoration(labelText: l10n.chatsTitleField),
        onSubmitted: (_) => _done(),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(onPressed: _done, child: Text(l10n.chatsCreate)),
      ],
    );
  }
}
