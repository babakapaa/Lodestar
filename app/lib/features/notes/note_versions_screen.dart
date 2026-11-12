import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note_actions.dart';
import 'notes_providers.dart';

class NoteVersionsScreen extends ConsumerWidget {
  const NoteVersionsScreen({required this.noteId, super.key});

  final String noteId;

  Future<void> _save(BuildContext context, WidgetRef ref) async {
    await NoteActions.remote(
      context,
      ref,
      () => ref.read(notesSyncProvider).saveVersion(noteId),
    );
    ref.invalidate(noteVersionsProvider(noteId));
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final note = ref.watch(noteProvider(noteId)).value;
    final versions = ref.watch(noteVersionsProvider(noteId));
    final list = versions.value;
    final format = DateFormat.yMMMd(l10n.localeName).add_Hm();
    return Scaffold(
      appBar: AppBar(title: Text(l10n.notesHistory)),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              Text(
                l10n.notesHistoryHint,
                style: TextStyle(color: context.palette.muted),
              ),
              const SizedBox(height: Spacing.m),
              if (note != null && note.role.canEdit && !note.deleted) ...[
                Align(
                  alignment: Alignment.centerLeft,
                  child: LodestarButton(
                    label: l10n.notesSaveVersion,
                    icon: Icons.bookmark_add_outlined,
                    variant: LodestarButtonVariant.secondary,
                    onPressed: () => _save(context, ref),
                  ),
                ),
                const SizedBox(height: Spacing.m),
              ],
              if (versions.hasError)
                FormError(message: ApiErrorText.of(l10n, versions.error!))
              else if (list == null)
                const Center(child: CircularProgressIndicator())
              else if (list.isEmpty)
                Text(l10n.notesHistoryEmpty)
              else
                for (final version in list.reversed)
                  ListTile(
                    key: ValueKey(version.id),
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.history),
                    title: Text(format.format(version.createdAt.toLocal())),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () =>
                        context.push(AppRoutes.noteVersion(noteId, version.id)),
                  ),
            ],
          ),
        ),
      ),
    );
  }
}
