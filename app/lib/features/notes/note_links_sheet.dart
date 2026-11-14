import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_actions.dart';
import 'note_link_tile.dart';
import 'note_resolver.dart';
import 'notes_providers.dart';

class NoteLinksSheet extends ConsumerWidget {
  const NoteLinksSheet({required this.note, super.key});

  final Note note;

  static Future<void> show(BuildContext context, Note note) =>
      showModalBottomSheet<void>(
        context: context,
        isScrollControlled: true,
        showDragHandle: true,
        builder: (context) => NoteLinksSheet(note: note),
      );

  void _open(BuildContext context, String id) {
    Navigator.of(context).pop();
    context.push(AppRoutes.note(id));
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final notes = ref.watch(notesProvider).value ?? const <Note>[];
    final backlinks = NoteResolver.backlinks(note, notes);
    final links = note.meta.links;
    return SafeArea(
      child: ListView(
        shrinkWrap: true,
        padding: const EdgeInsets.fromLTRB(Spacing.l, 0, Spacing.l, Spacing.l),
        children: [
          if (links.isEmpty && backlinks.isEmpty) Text(l10n.notesLinksNone),
          if (links.isNotEmpty)
            Section(
              title: l10n.notesLinksOut,
              child: Column(
                children: [
                  for (final link in links)
                    NoteLinkTile(
                      link: link,
                      target: NoteResolver.resolve(link, notes),
                      onOpen: (id) => _open(context, id),
                      onCreate: () {
                        Navigator.of(context).pop();
                        NoteActions.create(context, ref, path: link);
                      },
                    ),
                ],
              ),
            ),
          if (backlinks.isNotEmpty) ...[
            const SizedBox(height: Spacing.l),
            Section(
              title: l10n.notesLinksIn,
              child: Column(
                children: [
                  for (final other in backlinks)
                    ListTile(
                      contentPadding: EdgeInsets.zero,
                      leading: const Icon(Icons.call_received),
                      title: Text(other.meta.name),
                      subtitle: other.meta.folder.isEmpty
                          ? null
                          : Text(other.meta.folder),
                      onTap: () => _open(context, other.id),
                    ),
                ],
              ),
            ),
          ],
        ],
      ),
    );
  }
}
