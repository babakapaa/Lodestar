import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'notes_providers.dart';

class NoteVersionScreen extends ConsumerWidget {
  const NoteVersionScreen({
    required this.noteId,
    required this.versionId,
    super.key,
  });

  final String noteId;
  final String versionId;

  Future<void> _restore(
    BuildContext context,
    WidgetRef ref,
    String body,
  ) async {
    await ref.read(notesStoreProvider).replaceText(noteId, body);
    if (context.mounted) {
      context
        ..pop()
        ..pop();
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final note = ref.watch(noteProvider(noteId)).value;
    final body = ref.watch(noteVersionBodyProvider((noteId, versionId)));
    final text = body.value;
    final canRestore =
        text != null && note != null && note.role.canEdit && !note.deleted;
    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.notesVersion),
        actions: [
          if (canRestore)
            TextButton(
              onPressed: () => _restore(context, ref, text),
              child: Text(l10n.notesRestoreVersion),
            ),
        ],
      ),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 820),
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              if (body.hasError)
                FormError(message: ApiErrorText.of(l10n, body.error!))
              else if (text == null)
                const Center(child: CircularProgressIndicator())
              else
                SelectableText(text),
            ],
          ),
        ),
      ),
    );
  }
}
