import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_actions.dart';
import 'note_path_dialog.dart';
import 'note_role.dart';
import 'note_tags_dialog.dart';
import 'notes_providers.dart';

class NoteMenu extends ConsumerWidget {
  const NoteMenu({required this.note, required this.onDeleted, super.key});

  final Note note;
  final VoidCallback onDeleted;

  static Future<void> rename(
    BuildContext context,
    WidgetRef ref,
    Note note,
  ) async {
    final path = await NotePathDialog.show(context, note.meta.path);
    if (path != null && path != note.meta.path) {
      await ref
          .read(notesStoreProvider)
          .setMeta(note.id, note.meta.copyWith(path: path));
    }
  }

  Future<void> _tags(BuildContext context, WidgetRef ref) async {
    final tags = await NoteTagsDialog.show(context, note.meta.tags);
    if (tags != null) {
      await ref
          .read(notesStoreProvider)
          .setMeta(note.id, note.meta.copyWith(tags: tags));
    }
  }

  Future<void> _delete(BuildContext context, WidgetRef ref) async {
    final done = await NoteActions.remote(
      context,
      ref,
      () => ref.read(notesApiProvider).delete(note.id),
    );
    if (done) {
      onDeleted();
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final canEdit = note.role.canEdit && !note.deleted;
    final owner = note.role == NoteRole.owner;
    return PopupMenuButton<String>(
      onSelected: (action) => switch (action) {
        'rename' => rename(context, ref, note),
        'tags' => _tags(context, ref),
        'delete' => _delete(context, ref),
        'history' => context.push(AppRoutes.noteVersions(note.id)),
        _ => NoteActions.remote(
          context,
          ref,
          () => ref.read(notesApiProvider).restore(note.id),
        ),
      },
      itemBuilder: (context) => [
        if (canEdit) ...[
          PopupMenuItem(value: 'rename', child: Text(l10n.notesRename)),
          PopupMenuItem(value: 'tags', child: Text(l10n.notesTags)),
        ],
        PopupMenuItem(value: 'history', child: Text(l10n.notesHistory)),
        if (owner && !note.deleted)
          PopupMenuItem(value: 'delete', child: Text(l10n.notesDelete)),
        if (owner && note.deleted)
          PopupMenuItem(value: 'restore', child: Text(l10n.notesRestore)),
      ],
    );
  }
}
