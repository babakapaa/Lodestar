import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../core/navigation/layout_size.dart';
import '../../core/navigation/tab_page.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_actions.dart';
import 'note_tile.dart';
import 'notes_providers.dart';
import 'notes_sync_status.dart';
import 'notes_workspace.dart';

class NotesScreen extends ConsumerStatefulWidget {
  const NotesScreen({super.key});

  @override
  ConsumerState<NotesScreen> createState() => _NotesScreenState();
}

class _NotesScreenState extends ConsumerState<NotesScreen> {
  String _query = '';

  bool _matches(Note note) {
    final query = _query.trim().toLowerCase();
    if (query.isEmpty) {
      return true;
    }
    return note.meta.path.toLowerCase().contains(query) ||
        note.body.toLowerCase().contains(query) ||
        note.meta.tags.any((tag) => tag.toLowerCase().contains(query));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final status = ref.watch(notesSyncStatusProvider);
    final pending = ref.watch(notesPendingProvider).value ?? false;
    final all = ref.watch(notesProvider).value;
    final live = [
      for (final note in all ?? const <Note>[])
        if (!note.deleted) note,
    ];
    final shown = live.where(_matches).toList();
    final trashed = (all?.length ?? 0) - live.length;
    if (LayoutSize.of(MediaQuery.sizeOf(context).width) ==
        LayoutSize.expanded) {
      return TabPage(title: l10n.tabNotes, child: const NotesWorkspace());
    }
    return TabPage(
      title: l10n.tabNotes,
      actions: [
        IconButton(
          tooltip: l10n.notesGraph,
          icon: const Icon(Icons.hub_outlined),
          onPressed: () => context.push(AppRoutes.notesGraph),
        ),
        PopupMenuButton<String>(
          onSelected: (action) => switch (action) {
            'sources' => context.push(AppRoutes.sourcesPath),
            'export' => NoteActions.exportVault(ref),
            _ => NoteActions.importVault(context, ref),
          },
          itemBuilder: (context) => [
            PopupMenuItem(value: 'sources', child: Text(l10n.notesSources)),
            PopupMenuItem(value: 'export', child: Text(l10n.notesExport)),
            PopupMenuItem(value: 'import', child: Text(l10n.notesImport)),
          ],
        ),
      ],
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => NoteActions.create(context, ref),
        icon: const Icon(Icons.add),
        label: Text(l10n.notesNew),
      ),
      child: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: all == null
              ? const Center(child: CircularProgressIndicator())
              : RefreshIndicator(
                  onRefresh: ref.read(notesSyncStatusProvider.notifier).sync,
                  child: ListView(
                    padding: const EdgeInsets.fromLTRB(
                      Spacing.l,
                      Spacing.l,
                      Spacing.l,
                      96,
                    ),
                    children: [
                      if (status == NotesSyncStatus.offline && pending) ...[
                        NoticeBanner(
                          icon: Icons.cloud_off_outlined,
                          message: l10n.notesOffline,
                        ),
                        const SizedBox(height: Spacing.l),
                      ],
                      if (live.isEmpty)
                        EmptyState(
                          icon: Icons.sticky_note_2_outlined,
                          title: l10n.emptyNotesTitle,
                          body: l10n.emptyNotesBody,
                        )
                      else ...[
                        TextField(
                          decoration: InputDecoration(
                            prefixIcon: const Icon(Icons.search),
                            hintText: l10n.notesSearch,
                          ),
                          onChanged: (value) => setState(() => _query = value),
                        ),
                        const SizedBox(height: Spacing.s),
                        if (shown.isEmpty)
                          Padding(
                            padding: const EdgeInsets.all(Spacing.l),
                            child: Text(l10n.notesNothingFound),
                          ),
                        for (final note in shown)
                          NoteTile(key: ValueKey(note.id), note: note),
                      ],
                      if (trashed > 0)
                        ListTile(
                          contentPadding: EdgeInsets.zero,
                          leading: const Icon(Icons.delete_outline),
                          title: Text(l10n.notesTrash),
                          trailing: Text('$trashed'),
                          onTap: () => context.push(AppRoutes.notesTrash),
                        ),
                    ],
                  ),
                ),
        ),
      ),
    );
  }
}
