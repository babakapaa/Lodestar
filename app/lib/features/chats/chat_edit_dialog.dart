import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class ChatEditDialog extends StatefulWidget {
  const ChatEditDialog({required this.title, required this.about, super.key});

  final String title;
  final String about;

  static Future<(String, String)?> show(
    BuildContext context,
    String title,
    String about,
  ) => showDialog<(String, String)>(
    context: context,
    builder: (context) => ChatEditDialog(title: title, about: about),
  );

  @override
  State<ChatEditDialog> createState() => _ChatEditDialogState();
}

class _ChatEditDialogState extends State<ChatEditDialog> {
  late final _title = TextEditingController(text: widget.title);
  late final _about = TextEditingController(text: widget.about);

  @override
  void dispose() {
    _title.dispose();
    _about.dispose();
    super.dispose();
  }

  void _done() {
    final title = _title.text.trim();
    if (title.isNotEmpty) {
      Navigator.of(context).pop((title, _about.text.trim()));
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.chatsEditInfo),
      content: SizedBox(
        width: 480,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              key: const ValueKey('chat-edit-title'),
              controller: _title,
              autofocus: true,
              maxLength: 200,
              decoration: InputDecoration(labelText: l10n.chatsTitleLabel),
            ),
            TextField(
              key: const ValueKey('chat-edit-about'),
              controller: _about,
              minLines: 2,
              maxLines: 5,
              maxLength: 2000,
              decoration: InputDecoration(labelText: l10n.chatsAboutLabel),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(onPressed: _done, child: Text(l10n.actionSave)),
      ],
    );
  }
}
