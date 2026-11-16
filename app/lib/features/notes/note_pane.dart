import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'canvas/canvas_view.dart';
import 'note.dart';
import 'note_editor.dart';
import 'note_menu.dart';
import 'note_notices.dart';
import 'note_title_field.dart';

class NotePane extends ConsumerWidget {
  const NotePane({
    required this.note,
    required this.jumps,
    required this.onFollow,
    required this.onDeleted,
    super.key,
  });

  final Note note;
  final ValueListenable<int?> jumps;
  final ValueChanged<String> onFollow;
  final VoidCallback onDeleted;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final palette = context.palette;
    final canEdit = note.role.canEdit && !note.deleted;
    final textTheme = Theme.of(context).textTheme;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        SizedBox(
          height: 44,
          child: Row(
            children: [
              const SizedBox(width: Spacing.l),
              if (note.meta.folder.isNotEmpty) ...[
                Text(
                  note.meta.folder.replaceAll('/', ' / '),
                  style: textTheme.bodyMedium!.copyWith(color: palette.muted),
                ),
                Text(
                  ' / ',
                  style: textTheme.bodyMedium!.copyWith(color: palette.muted),
                ),
              ],
              Flexible(
                child: Text(
                  note.meta.label,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: textTheme.titleSmall,
                ),
              ),
              const Spacer(),
              if (!note.deleted)
                IconButton(
                  tooltip: l10n.notesShare,
                  icon: const Icon(Icons.group_outlined),
                  onPressed: () => context.push(AppRoutes.noteMembers(note.id)),
                ),
              NoteMenu(note: note, onDeleted: onDeleted),
              const SizedBox(width: Spacing.xs),
            ],
          ),
        ),
        const Divider(),
        NoteNotices(note: note),
        Expanded(
          child: note.meta.isCanvas
              ? CanvasView(
                  key: ValueKey(note.id),
                  noteId: note.id,
                  readOnly: !canEdit,
                  onOpenNote: onFollow,
                )
              : Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    if (!note.headedByTitle)
                      Align(
                        alignment: Alignment.topCenter,
                        child: ConstrainedBox(
                          constraints: const BoxConstraints(
                            maxWidth: NoteEditor.readableWidth,
                          ),
                          child: NoteTitleField(
                            key: ValueKey('title-${note.id}'),
                            note: note,
                            readOnly: !canEdit,
                          ),
                        ),
                      ),
                    Expanded(
                      child: NoteEditor(
                        key: ValueKey(note.id),
                        noteId: note.id,
                        readOnly: !canEdit,
                        jumps: jumps,
                        onOpenLink: onFollow,
                      ),
                    ),
                  ],
                ),
        ),
      ],
    );
  }
}
