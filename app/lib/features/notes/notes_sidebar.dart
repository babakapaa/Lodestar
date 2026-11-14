import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/notice_banner.dart';
import '../../design/components/pane_heading.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_folder.dart';
import 'note_resolver.dart';
import 'note_tree.dart';
import 'note_tree_row.dart';
import 'notes_providers.dart';
import 'notes_sync_status.dart';

class NotesSidebar extends ConsumerStatefulWidget {
  const NotesSidebar({
    required this.active,
    required this.onOpen,
    required this.onCreate,
    required this.onQuick,
    required this.onGraph,
    super.key,
  });

  final String? active;
  final void Function(Note note) onOpen;
  final VoidCallback onCreate;
  final VoidCallback onQuick;
  final VoidCallback onGraph;

  @override
  ConsumerState<NotesSidebar> createState() => _NotesSidebarState();
}

class _NotesSidebarState extends ConsumerState<NotesSidebar> {
  String _query = '';

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final all = ref.watch(notesProvider).value ?? const <Note>[];
    final status = ref.watch(notesSyncStatusProvider);
    final pending = ref.watch(notesPendingProvider).value ?? false;
    final trashed = all.where((note) => note.deleted).length;
    final needle = _query.trim().toLowerCase();
    final shown = needle.isEmpty
        ? null
        : [
            for (final note in NoteResolver.search('', all))
              if (note.meta.path.toLowerCase().contains(needle) ||
                  note.body.toLowerCase().contains(needle) ||
                  note.meta.tags.any(
                    (tag) => tag.toLowerCase().contains(needle),
                  ))
                note,
          ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        PaneHeading(
          title: l10n.notesVault,
          trailing: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              IconButton(
                tooltip: l10n.notesNew,
                visualDensity: VisualDensity.compact,
                icon: const Icon(Icons.edit_note),
                onPressed: widget.onCreate,
              ),
              IconButton(
                tooltip: l10n.notesQuickOpen,
                visualDensity: VisualDensity.compact,
                icon: const Icon(Icons.manage_search),
                onPressed: widget.onQuick,
              ),
              IconButton(
                tooltip: l10n.notesGraph,
                visualDensity: VisualDensity.compact,
                icon: const Icon(Icons.hub_outlined),
                onPressed: widget.onGraph,
              ),
            ],
          ),
        ),
        Padding(
          padding: const EdgeInsets.fromLTRB(
            Spacing.m,
            Spacing.xs,
            Spacing.m,
            Spacing.s,
          ),
          child: TextField(
            decoration: InputDecoration(
              isDense: true,
              prefixIcon: const Icon(Icons.search, size: 18),
              hintText: l10n.notesSearch,
            ),
            onChanged: (value) => setState(() => _query = value),
          ),
        ),
        if (status == NotesSyncStatus.offline && pending)
          Padding(
            padding: const EdgeInsets.all(Spacing.m),
            child: NoticeBanner(
              icon: Icons.cloud_off_outlined,
              message: l10n.notesOffline,
            ),
          ),
        Expanded(
          child: shown == null
              ? NoteTree(
                  root: NoteFolder.of(all),
                  active: widget.active,
                  onOpen: widget.onOpen,
                )
              : shown.isEmpty
              ? Padding(
                  padding: const EdgeInsets.all(Spacing.l),
                  child: Text(l10n.notesNothingFound),
                )
              : ListView(
                  children: [
                    for (final note in shown)
                      NoteTreeRow(
                        key: ValueKey('found:${note.id}'),
                        depth: 0,
                        icon: Icons.description_outlined,
                        label: note.meta.path,
                        selected: note.id == widget.active,
                        onTap: () => widget.onOpen(note),
                      ),
                  ],
                ),
        ),
        const Divider(),
        const SizedBox(height: Spacing.xs),
        NoteTreeRow(
          depth: 0,
          icon: Icons.menu_book_outlined,
          label: l10n.notesSources,
          onTap: () => context.push(AppRoutes.sourcesPath),
        ),
        NoteTreeRow(
          depth: 0,
          icon: Icons.delete_outline,
          label: trashed > 0
              ? '${l10n.notesTrash} ($trashed)'
              : l10n.notesTrash,
          onTap: () => context.push(AppRoutes.notesTrash),
        ),
        const SizedBox(height: Spacing.s),
      ],
    );
  }
}
