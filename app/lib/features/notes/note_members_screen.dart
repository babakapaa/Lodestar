import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../people/person_picker_dialog.dart';
import 'note_access_dialog.dart';
import 'note_actions.dart';
import 'note_member_tile.dart';
import 'note_role.dart';
import 'notes_providers.dart';

class NoteMembersScreen extends ConsumerWidget {
  const NoteMembersScreen({required this.noteId, super.key});

  final String noteId;

  Future<void> _add(BuildContext context, WidgetRef ref) async {
    final person = await PersonPickerDialog.show(context);
    if (person == null || !context.mounted) {
      return;
    }
    final canEdit = await NoteAccessDialog.show(context, person.name);
    if (canEdit == null || !context.mounted) {
      return;
    }
    await _share(context, ref, person.id, canEdit: canEdit);
  }

  Future<void> _share(
    BuildContext context,
    WidgetRef ref,
    String userId, {
    required bool canEdit,
  }) async {
    await NoteActions.remote(
      context,
      ref,
      () => ref.read(notesApiProvider).share(noteId, userId, canEdit: canEdit),
    );
    ref.invalidate(noteMembersProvider(noteId));
  }

  Future<void> _remove(
    BuildContext context,
    WidgetRef ref,
    String userId,
  ) async {
    await NoteActions.remote(
      context,
      ref,
      () => ref.read(notesApiProvider).unshare(noteId, userId),
    );
    ref.invalidate(noteMembersProvider(noteId));
  }

  Future<void> _leave(BuildContext context, WidgetRef ref) async {
    final left = await NoteActions.remote(context, ref, () async {
      await ref.read(notesApiProvider).leave(noteId);
      await ref.read(notesStoreProvider).forget(noteId);
    });
    if (left && context.mounted) {
      context.go(AppRoutes.notesPath);
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final note = ref.watch(noteProvider(noteId)).value;
    final owner = note?.role == NoteRole.owner;
    final members = ref.watch(noteMembersProvider(noteId));
    final list = members.value;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.notesShare)),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              if (note != null)
                Text(
                  note.meta.name,
                  style: Theme.of(context).textTheme.titleMedium,
                ),
              const SizedBox(height: Spacing.m),
              if (members.hasError)
                FormError(message: ApiErrorText.of(l10n, members.error!))
              else if (list == null)
                const Center(child: CircularProgressIndicator())
              else if (list.isEmpty)
                Text(l10n.notesShareEmpty)
              else
                for (final member in list)
                  NoteMemberTile(
                    key: ValueKey(member.userId),
                    member: member,
                    manage: owner,
                    onAccess: (canEdit) =>
                        _share(context, ref, member.userId, canEdit: canEdit),
                    onRemove: () => _remove(context, ref, member.userId),
                  ),
              const SizedBox(height: Spacing.l),
              Align(
                alignment: Alignment.centerLeft,
                child: owner
                    ? LodestarButton(
                        label: l10n.notesShareAdd,
                        icon: Icons.person_add_alt,
                        onPressed: () => _add(context, ref),
                      )
                    : LodestarButton(
                        label: l10n.notesLeave,
                        icon: Icons.logout,
                        variant: LodestarButtonVariant.secondary,
                        onPressed: () => _leave(context, ref),
                      ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
