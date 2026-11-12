import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_actions.dart';
import 'note_editor.dart';
import 'note_links_sheet.dart';
import 'note_path_dialog.dart';
import 'note_role.dart';
import 'note_tags_dialog.dart';
import 'notes_providers.dart';

class NoteScreen extends ConsumerWidget {
  const NoteScreen({required this.noteId, super.key});

  final String noteId;

  Future<void> _rename(BuildContext context, WidgetRef ref, Note note) async {
    final path = await NotePathDialog.show(context, note.meta.path);
    if (path != null && path != note.meta.path) {
      await ref
          .read(notesStoreProvider)
          .setMeta(noteId, note.meta.copyWith(path: path));
    }
  }

  Future<void> _tags(BuildContext context, WidgetRef ref, Note note) async {
    final tags = await NoteTagsDialog.show(context, note.meta.tags);
    if (tags != null) {
      await ref
          .read(notesStoreProvider)
          .setMeta(noteId, note.meta.copyWith(tags: tags));
    }
  }

  Future<void> _delete(BuildContext context, WidgetRef ref) async {
    final done = await NoteActions.remote(
      context,
      ref,
      () => ref.read(notesApiProvider).delete(noteId),
    );
    if (done && context.mounted) {
      context.pop();
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final value = ref.watch(noteProvider(noteId));
    final note = value.value;
    if (note == null) {
      return Scaffold(
        appBar: AppBar(),
        body: value.isLoading
            ? const Center(child: CircularProgressIndicator())
            : EmptyState(
                icon: Icons.search_off,
                title: l10n.notesNotFound,
                body: l10n.errorNoteNotFound,
              ),
      );
    }
    final canEdit = note.role.canEdit && !note.deleted;
    final owner = note.role == NoteRole.owner;
    final error = note.error;
    return Scaffold(
      appBar: AppBar(
        title: GestureDetector(
          onTap: canEdit ? () => _rename(context, ref, note) : null,
          child: Text(note.meta.name),
        ),
        actions: [
          IconButton(
            tooltip: l10n.notesLinks,
            icon: const Icon(Icons.hub_outlined),
            onPressed: () => NoteLinksSheet.show(context, note),
          ),
          if (!note.deleted)
            IconButton(
              tooltip: l10n.notesShare,
              icon: const Icon(Icons.group_outlined),
              onPressed: () => context.push(AppRoutes.noteMembers(noteId)),
            ),
          PopupMenuButton<String>(
            onSelected: (action) => switch (action) {
              'rename' => _rename(context, ref, note),
              'tags' => _tags(context, ref, note),
              'delete' => _delete(context, ref),
              'history' => context.push(AppRoutes.noteVersions(noteId)),
              _ => NoteActions.remote(
                context,
                ref,
                () => ref.read(notesApiProvider).restore(noteId),
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
          ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 820),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                for (final message in [
                  if (note.deleted) l10n.notesInTrash,
                  if (!note.role.canEdit) l10n.notesReadOnly,
                  if (error != null)
                    l10n.notesNotAccepted(ApiErrorText.ofCode(l10n, error)),
                ])
                  Padding(
                    padding: const EdgeInsets.fromLTRB(
                      Spacing.l,
                      Spacing.m,
                      Spacing.l,
                      0,
                    ),
                    child: NoticeBanner(
                      icon: Icons.info_outline,
                      message: message,
                    ),
                  ),
                Expanded(
                  child: NoteEditor(
                    key: ValueKey(noteId),
                    noteId: noteId,
                    readOnly: !canEdit,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
