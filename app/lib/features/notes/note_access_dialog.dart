import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class NoteAccessDialog extends StatelessWidget {
  const NoteAccessDialog({required this.name, super.key});

  final String name;

  static Future<bool?> show(BuildContext context, String name) =>
      showDialog<bool>(
        context: context,
        builder: (context) => NoteAccessDialog(name: name),
      );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return SimpleDialog(
      title: Text(name),
      children: [
        SimpleDialogOption(
          onPressed: () => Navigator.of(context).pop(true),
          child: Text(l10n.notesCanEdit),
        ),
        SimpleDialogOption(
          onPressed: () => Navigator.of(context).pop(false),
          child: Text(l10n.notesCanRead),
        ),
      ],
    );
  }
}
