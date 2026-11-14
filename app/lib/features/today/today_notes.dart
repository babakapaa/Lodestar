import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../core/navigation/layout_size.dart';
import '../../design/components/section.dart';
import '../../design/components/tile_group.dart';
import '../../l10n/l10n_context.dart';
import '../notes/note.dart';
import '../notes/notes_providers.dart';
import 'today_placeholder.dart';

class TodayNotes extends ConsumerWidget {
  const TodayNotes({super.key});

  static const limit = 5;

  void _open(BuildContext context, WidgetRef ref, Note note) {
    if (LayoutSize.of(MediaQuery.sizeOf(context).width) ==
        LayoutSize.expanded) {
      ref.read(notesTabsProvider.notifier).open(note.id);
      context.go(AppRoutes.notesPath);
    } else {
      context.push(AppRoutes.note(note.id));
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final notes = [
      ...?ref.watch(notesProvider).value?.where((note) => !note.deleted),
    ]..sort((a, b) => b.changedAt.compareTo(a.changedAt));
    return Section(
      title: l10n.todayNotes,
      child: notes.isEmpty
          ? TodayPlaceholder(
              icon: Icons.sticky_note_2_outlined,
              text: l10n.todayNoNotes,
            )
          : TileGroup(
              children: [
                for (final note in notes.take(limit))
                  ListTile(
                    key: ValueKey('today-note:${note.id}'),
                    contentPadding: TileGroup.tilePadding,
                    leading: Icon(
                      note.meta.isCanvas
                          ? Icons.dashboard_outlined
                          : Icons.description_outlined,
                    ),
                    title: Text(
                      note.meta.label,
                      overflow: TextOverflow.ellipsis,
                    ),
                    subtitle: note.meta.folder.isEmpty
                        ? null
                        : Text(note.meta.folder),
                    onTap: () => _open(context, ref, note),
                  ),
              ],
            ),
    );
  }
}
