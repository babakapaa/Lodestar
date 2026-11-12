import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'note.dart';

class NoteLinkTile extends StatelessWidget {
  const NoteLinkTile({
    super.key,
    required this.link,
    required this.target,
    required this.onOpen,
    required this.onCreate,
  });

  final String link;
  final Note? target;
  final void Function(String id) onOpen;
  final VoidCallback onCreate;

  @override
  Widget build(BuildContext context) {
    final note = target;
    if (note == null) {
      return ListTile(
        contentPadding: EdgeInsets.zero,
        leading: const Icon(Icons.add_link),
        title: Text(context.l10n.notesLinkCreate(link)),
        onTap: onCreate,
      );
    }
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: const Icon(Icons.north_east),
      title: Text(note.meta.name),
      subtitle: note.meta.folder.isEmpty ? null : Text(note.meta.folder),
      onTap: () => onOpen(note.id),
    );
  }
}
