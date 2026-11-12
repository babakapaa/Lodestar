import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_actions.dart';
import 'note_role.dart';
import 'note_tile.dart';
import 'notes_providers.dart';

class NotesTrashScreen extends ConsumerWidget {
  const NotesTrashScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final deleted = [
      for (final note in ref.watch(notesProvider).value ?? const <Note>[])
        if (note.deleted) note,
    ];
    return Scaffold(
      appBar: AppBar(title: Text(l10n.notesTrash)),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              Text(
                l10n.notesTrashHint,
                style: TextStyle(color: context.palette.muted),
              ),
              const SizedBox(height: Spacing.m),
              if (deleted.isEmpty) Text(l10n.notesTrashEmpty),
              for (final note in deleted)
                NoteTile(
                  key: ValueKey(note.id),
                  note: note,
                  trailing: note.role == NoteRole.owner
                      ? TextButton(
                          onPressed: () => NoteActions.remote(
                            context,
                            ref,
                            () => ref.read(notesApiProvider).restore(note.id),
                          ),
                          child: Text(l10n.notesRestore),
                        )
                      : null,
                ),
            ],
          ),
        ),
      ),
    );
  }
}
