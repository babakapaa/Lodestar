import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'note_member.dart';

class NoteMemberTile extends StatelessWidget {
  const NoteMemberTile({
    required this.member,
    required this.manage,
    required this.onAccess,
    required this.onRemove,
    super.key,
  });

  final NoteMember member;
  final bool manage;
  final void Function(bool canEdit) onAccess;
  final VoidCallback onRemove;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: const Icon(Icons.person_outline),
      title: Text(member.label),
      subtitle: Text(
        '${member.email}\n${member.canEdit ? l10n.notesCanEdit : l10n.notesCanRead}',
      ),
      isThreeLine: true,
      trailing: manage
          ? PopupMenuButton<String>(
              onSelected: (action) => switch (action) {
                'edit' => onAccess(true),
                'read' => onAccess(false),
                _ => onRemove(),
              },
              itemBuilder: (context) => [
                PopupMenuItem(value: 'edit', child: Text(l10n.notesCanEdit)),
                PopupMenuItem(value: 'read', child: Text(l10n.notesCanRead)),
                PopupMenuItem(value: 'remove', child: Text(l10n.notesUnshare)),
              ],
            )
          : null,
    );
  }
}
