import 'dart:async';

import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/empty_state.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_actions.dart';
import 'note_inspector.dart';
import 'note_pane.dart';
import 'notes_graph_view.dart';
import 'notes_providers.dart';
import 'notes_sidebar.dart';
import 'notes_tab_bar.dart';
import 'notes_tabs.dart';
import 'quick_switcher.dart';

class NotesWorkspace extends ConsumerStatefulWidget {
  const NotesWorkspace({super.key});

  static const sidebarWidth = 272.0;
  static const inspectorWidth = 280.0;

  @override
  ConsumerState<NotesWorkspace> createState() => _NotesWorkspaceState();
}

class _NotesWorkspaceState extends ConsumerState<NotesWorkspace> {
  final _jumps = ValueNotifier<int?>(null);
  bool _inspector = true;

  void _open(String id, {bool? newTab}) {
    final keyboard = HardwareKeyboard.instance;
    ref
        .read(notesTabsProvider.notifier)
        .open(
          id,
          newTab:
              newTab ?? (keyboard.isControlPressed || keyboard.isMetaPressed),
        );
  }

  Future<void> _create() async {
    final id = await NoteActions.add(context, ref);
    _open(id, newTab: true);
  }

  Future<void> _createCanvas() async {
    final id = await NoteActions.addCanvas(context, ref);
    _open(id, newTab: true);
  }

  Future<void> _quick() async {
    final notes = ref.read(notesProvider).value ?? const <Note>[];
    final choice = await QuickSwitcher.show(context, notes);
    if (choice == null || !mounted) {
      return;
    }
    final id = choice.noteId;
    if (id != null) {
      _open(id, newTab: false);
      return;
    }
    await _follow(choice.createPath!);
  }

  Future<void> _follow(String link) =>
      NoteActions.follow(context, ref, link, (id) => _open(id, newTab: false));

  void _jump(int offset) {
    _jumps.value = null;
    _jumps.value = offset;
  }

  void _closeActive() {
    final active = ref.read(notesTabsProvider).active;
    if (active != null) {
      ref.read(notesTabsProvider.notifier).close(active);
    }
  }

  Map<ShortcutActivator, VoidCallback> get _shortcuts => {
    for (final control in [true, false]) ...{
      SingleActivator(
        LogicalKeyboardKey.keyO,
        control: control,
        meta: !control,
      ): () =>
          unawaited(_quick()),
      SingleActivator(
        LogicalKeyboardKey.keyN,
        control: control,
        meta: !control,
      ): () =>
          unawaited(_create()),
      SingleActivator(
        LogicalKeyboardKey.keyW,
        control: control,
        meta: !control,
      ): _closeActive,
      SingleActivator(
        LogicalKeyboardKey.keyG,
        control: control,
        meta: !control,
      ): () =>
          _open(NotesTabs.graph, newTab: true),
    },
  };

  @override
  void dispose() {
    _jumps.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = context.palette;
    final tabs = ref.watch(notesTabsProvider);
    final notes = ref.watch(notesProvider).value;
    final active = tabs.active;
    Note? note;
    for (final candidate in notes ?? const <Note>[]) {
      if (candidate.id == active) {
        note = candidate;
      }
    }
    final Widget content;
    if (notes == null) {
      content = const Center(child: CircularProgressIndicator());
    } else if (active == NotesTabs.graph) {
      content = NotesGraphView(onOpen: (id) => _open(id, newTab: false));
    } else if (note != null) {
      content = NotePane(
        key: ValueKey('pane:${note.id}'),
        note: note,
        jumps: _jumps,
        onFollow: _follow,
        onDeleted: () => ref.read(notesTabsProvider.notifier).close(note!.id),
      );
    } else {
      content = Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          EmptyState(
            icon: Icons.sticky_note_2_outlined,
            title: notes.any((n) => !n.deleted)
                ? l10n.notesNoTabsTitle
                : l10n.emptyNotesTitle,
            body: notes.any((n) => !n.deleted)
                ? l10n.notesNoTabsBody
                : l10n.emptyNotesBody,
          ),
          const SizedBox(height: Spacing.s),
          Wrap(
            spacing: Spacing.s,
            children: [
              FilledButton.icon(
                onPressed: _create,
                icon: const Icon(Icons.add),
                label: Text(l10n.notesNew),
              ),
              OutlinedButton.icon(
                onPressed: _quick,
                icon: const Icon(Icons.manage_search),
                label: Text(l10n.notesQuickOpen),
              ),
            ],
          ),
        ],
      );
    }
    final showInspector = _inspector && note != null && !note.meta.isCanvas;
    return CallbackShortcuts(
      bindings: _shortcuts,
      child: Focus(
        autofocus: true,
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Container(
              width: NotesWorkspace.sidebarWidth,
              color: palette.side,
              child: NotesSidebar(
                active: active,
                onOpen: (note) => _open(note.id),
                onCreate: _create,
                onCanvas: _createCanvas,
                onQuick: _quick,
                onGraph: () => _open(NotesTabs.graph, newTab: true),
              ),
            ),
            const VerticalDivider(width: 1),
            Expanded(
              child: ColoredBox(
                color: palette.surface,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    NotesTabBar(
                      trailing: [
                        IconButton(
                          tooltip: l10n.notesNew,
                          icon: const Icon(Icons.add),
                          onPressed: _create,
                        ),
                        IconButton(
                          tooltip: l10n.notesInspector,
                          isSelected: _inspector,
                          icon: const Icon(Icons.view_sidebar_outlined),
                          onPressed: () =>
                              setState(() => _inspector = !_inspector),
                        ),
                      ],
                    ),
                    const Divider(),
                    Expanded(child: content),
                  ],
                ),
              ),
            ),
            if (showInspector) ...[
              const VerticalDivider(width: 1),
              Container(
                width: NotesWorkspace.inspectorWidth,
                color: palette.side,
                child: NoteInspector(
                  note: note,
                  onOpen: (id) => _open(id, newTab: false),
                  onFollow: _follow,
                  onJump: _jump,
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
