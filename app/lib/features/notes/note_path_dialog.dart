import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'note_meta.dart';

class NotePathDialog extends StatefulWidget {
  const NotePathDialog({required this.path, super.key});

  final String path;

  static Future<String?> show(BuildContext context, String path) =>
      showDialog<String>(
        context: context,
        builder: (context) => NotePathDialog(path: path),
      );

  @override
  State<NotePathDialog> createState() => _NotePathDialogState();
}

class _NotePathDialogState extends State<NotePathDialog> {
  late final _controller = TextEditingController(text: widget.path);

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _submit() {
    final path = NoteMeta.normalizePath(_controller.text);
    if (path.isNotEmpty) {
      Navigator.of(context).pop(path);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.notesRename),
      content: TextField(
        controller: _controller,
        autofocus: true,
        decoration: InputDecoration(labelText: l10n.notesPathLabel),
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
