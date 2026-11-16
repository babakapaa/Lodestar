import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_actions.dart';
import 'note_folder.dart';
import 'note_meta.dart';
import 'note_role.dart';
import 'note_tree_row.dart';
import 'notes_providers.dart';

class NoteTree extends ConsumerStatefulWidget {
  const NoteTree({
    required this.root,
    required this.active,
    required this.onOpen,
    required this.onCreate,
    super.key,
  });

  final NoteFolder root;
  final String? active;
  final void Function(Note note) onOpen;
  final ValueChanged<String> onCreate;

  @override
  ConsumerState<NoteTree> createState() => _NoteTreeState();
}

class _NoteTreeState extends ConsumerState<NoteTree> {
  final _collapsed = <String>{};
  String? _editing;
  String? _target;

  void _toggle(String path) => setState(() {
    if (!_collapsed.remove(path)) {
      _collapsed.add(path);
    }
  });

  Future<void> _menu(
    Offset at,
    List<(String, IconData, VoidCallback)> items,
  ) async {
    final overlay =
        Overlay.of(context).context.findRenderObject()! as RenderBox;
    final picked = await showMenu<int>(
      context: context,
      position: RelativeRect.fromRect(
        at & const Size(1, 1),
        Offset.zero & overlay.size,
      ),
      items: [
        for (var i = 0; i < items.length; i++)
          PopupMenuItem(
            value: i,
            height: 36,
            child: Row(
              children: [
                Icon(items[i].$2, size: 18),
                const SizedBox(width: Spacing.m),
                Text(items[i].$1),
              ],
            ),
          ),
      ],
    );
    if (picked != null) {
      items[picked].$3();
    }
  }

  void _folderMenu(String folder, Offset at) {
    final l10n = context.l10n;
    _menu(at, [
      (l10n.notesNewHere, Icons.edit_note, () => widget.onCreate(folder)),
      (
        l10n.notesNewFolder,
        Icons.create_new_folder_outlined,
        () => widget.onCreate(
          folder.isEmpty
              ? l10n.notesNewFolder
              : '$folder/${l10n.notesNewFolder}',
        ),
      ),
      if (folder.isNotEmpty)
        (
          l10n.notesRenameShort,
          Icons.drive_file_rename_outline,
          () => setState(() => _editing = 'folder:$folder'),
        ),
    ]);
  }

  void _noteMenu(Note note, Offset at) {
    final l10n = context.l10n;
    _menu(at, [
      if (note.role.canEdit)
        (
          l10n.notesRenameShort,
          Icons.drive_file_rename_outline,
          () => setState(() => _editing = 'note:${note.id}'),
        ),
      if (note.role == NoteRole.owner)
        (
          l10n.notesDelete,
          Icons.delete_outline,
          () => NoteActions.remote(
            context,
            ref,
            () => ref.read(notesApiProvider).delete(note.id),
          ),
        ),
    ]);
  }

  Future<void> _renameNote(Note note, String? name) async {
    setState(() => _editing = null);
    final clean = name?.replaceAll('/', ' ').trim() ?? '';
    if (clean.isEmpty || clean == note.meta.label) {
      return;
    }
    final meta = note.meta;
    final file = meta.isCanvas ? '$clean${NoteMeta.canvasExtension}' : clean;
    await ref
        .read(notesStoreProvider)
        .setMeta(
          note.id,
          meta.copyWith(
            path: NoteMeta.normalizePath(
              meta.folder.isEmpty ? file : '${meta.folder}/$file',
            ),
          ),
        );
  }

  Future<void> _renameFolder(NoteFolder folder, String? name) async {
    setState(() => _editing = null);
    final clean = name?.replaceAll('/', ' ').trim() ?? '';
    if (clean.isEmpty || clean == folder.name) {
      return;
    }
    final at = folder.path.lastIndexOf('/');
    final parent = at < 0 ? '' : folder.path.substring(0, at);
    await NoteActions.renameFolder(
      ref,
      folder.path,
      parent.isEmpty ? clean : '$parent/$clean',
    );
  }

  Widget _drop(String folder, Widget child) {
    return DragTarget<Note>(
      onWillAcceptWithDetails: (details) {
        final ok =
            details.data.meta.folder != folder && details.data.role.canEdit;
        if (ok) {
          setState(() => _target = folder);
        }
        return ok;
      },
      onLeave: (_) => setState(() => _target = null),
      onAcceptWithDetails: (details) {
        setState(() => _target = null);
        NoteActions.move(ref, details.data, folder);
      },
      builder: (context, _, _) => child,
    );
  }

  void _rows(NoteFolder folder, int depth, List<Widget> rows) {
    for (final child in folder.folders) {
      final open = !_collapsed.contains(child.path);
      rows.add(
        _drop(
          child.path,
          NoteTreeRow(
            key: ValueKey('folder:${child.path}'),
            depth: depth,
            icon: open ? Icons.folder_open_outlined : Icons.folder_outlined,
            label: child.name,
            strong: true,
            hovered: _target == child.path,
            editing: _editing == 'folder:${child.path}',
            onRename: (name) => _renameFolder(child, name),
            onMenu: (at) => _folderMenu(child.path, at),
            onTap: () => _toggle(child.path),
          ),
        ),
      );
      if (open) {
        _rows(child, depth + 1, rows);
      }
    }
    for (final note in folder.notes) {
      final row = NoteTreeRow(
        key: ValueKey('note:${note.id}'),
        depth: depth,
        icon: note.error != null
            ? Icons.sync_problem
            : note.meta.isCanvas
            ? Icons.dashboard_outlined
            : Icons.description_outlined,
        label: note.meta.label,
        selected: note.id == widget.active,
        editing: _editing == 'note:${note.id}',
        onRename: (name) => _renameNote(note, name),
        onMenu: (at) => _noteMenu(note, at),
        onTap: () => widget.onOpen(note),
      );
      rows.add(
        note.role.canEdit
            ? Draggable<Note>(
                data: note,
                feedback: _DragChip(label: note.meta.label),
                childWhenDragging: Opacity(opacity: 0.4, child: row),
                child: row,
              )
            : row,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final rows = <Widget>[];
    _rows(widget.root, 0, rows);
    return _drop(
      '',
      GestureDetector(
        behavior: HitTestBehavior.translucent,
        onSecondaryTapDown: (details) =>
            _folderMenu('', details.globalPosition),
        child: ColoredBox(
          color: _target == ''
              ? context.palette.accent.withValues(alpha: 0.06)
              : Colors.transparent,
          child: ListView(children: rows),
        ),
      ),
    );
  }
}

class _DragChip extends StatelessWidget {
  const _DragChip({required this.label});

  final String label;

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    return Material(
      color: palette.surface,
      elevation: 4,
      borderRadius: context.tokens.borderRadiusSmall,
      child: Padding(
        padding: const EdgeInsets.symmetric(
          horizontal: Spacing.m,
          vertical: Spacing.s,
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.description_outlined, size: 16, color: palette.accent),
            const SizedBox(width: Spacing.s),
            Text(label, style: Theme.of(context).textTheme.bodyMedium),
          ],
        ),
      ),
    );
  }
}
