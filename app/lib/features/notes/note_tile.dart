import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_role.dart';

class NoteTile extends StatelessWidget {
  const NoteTile({required this.note, this.trailing, super.key});

  final Note note;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final preview = note.body.trim().split('\n').first;
    final details = [
      if (note.meta.folder.isNotEmpty) note.meta.folder,
      if (note.role != NoteRole.owner) l10n.notesSharedNote,
      for (final tag in note.meta.tags) '#$tag',
    ].join(', ');
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Icon(
        note.error != null ? Icons.sync_problem : Icons.description_outlined,
        color: note.error != null ? context.palette.warn : null,
      ),
      title: Text(note.meta.name),
      subtitle: Text(
        [
          if (details.isNotEmpty) details,
          if (preview.isNotEmpty) preview,
        ].join('\n'),
        maxLines: 2,
        overflow: TextOverflow.ellipsis,
      ),
      trailing: trailing ?? const Icon(Icons.chevron_right),
      onTap: () => context.push(AppRoutes.note(note.id)),
    );
  }
}
