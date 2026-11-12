import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class NoteTagsDialog extends StatefulWidget {
  const NoteTagsDialog({required this.tags, super.key});

  final List<String> tags;

  static Future<List<String>?> show(BuildContext context, List<String> tags) =>
      showDialog<List<String>>(
        context: context,
        builder: (context) => NoteTagsDialog(tags: tags),
      );

  @override
  State<NoteTagsDialog> createState() => _NoteTagsDialogState();
}

class _NoteTagsDialogState extends State<NoteTagsDialog> {
  late final _controller = TextEditingController(text: widget.tags.join(', '));

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _submit() {
    final tags = <String>{
      for (final part in _controller.text.split(','))
        if (part.trim().replaceAll('#', '').isNotEmpty)
          part.trim().replaceAll('#', ''),
    };
    Navigator.of(context).pop(tags.toList());
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.notesTags),
      content: TextField(
        controller: _controller,
        autofocus: true,
        decoration: InputDecoration(labelText: l10n.notesTagsLabel),
        onSubmitted: (_) => _submit(),
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
