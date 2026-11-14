import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/empty_state.dart';
import '../../l10n/l10n_context.dart';
import 'canvas/canvas_view.dart';
import 'note_actions.dart';
import 'note_editor.dart';
import 'note_links_sheet.dart';
import 'note_menu.dart';
import 'note_notices.dart';
import 'notes_providers.dart';

class NoteScreen extends ConsumerWidget {
  const NoteScreen({required this.noteId, super.key});

  final String noteId;

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
    void follow(String link) => NoteActions.follow(
      context,
      ref,
      link,
      (id) => context.push(AppRoutes.note(id)),
    );
    return Scaffold(
      appBar: AppBar(
        title: GestureDetector(
          onTap: canEdit ? () => NoteMenu.rename(context, ref, note) : null,
          child: Text(note.meta.label),
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
          NoteMenu(note: note, onDeleted: () => context.pop()),
        ],
      ),
      body: SafeArea(
        top: false,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            NoteNotices(note: note),
            Expanded(
              child: note.meta.isCanvas
                  ? CanvasView(
                      key: ValueKey(noteId),
                      noteId: noteId,
                      readOnly: !canEdit,
                      onOpenNote: follow,
                    )
                  : NoteEditor(
                      key: ValueKey(noteId),
                      noteId: noteId,
                      readOnly: !canEdit,
                      onOpenLink: follow,
                    ),
            ),
          ],
        ),
      ),
    );
  }
}
