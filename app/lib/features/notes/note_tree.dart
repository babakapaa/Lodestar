import 'package:material_ui/material_ui.dart';

import 'note.dart';
import 'note_folder.dart';
import 'note_tree_row.dart';

class NoteTree extends StatefulWidget {
  const NoteTree({
    required this.root,
    required this.active,
    required this.onOpen,
    super.key,
  });

  final NoteFolder root;
  final String? active;
  final void Function(Note note) onOpen;

  @override
  State<NoteTree> createState() => _NoteTreeState();
}

class _NoteTreeState extends State<NoteTree> {
  final _collapsed = <String>{};

  void _toggle(String path) => setState(() {
    if (!_collapsed.remove(path)) {
      _collapsed.add(path);
    }
  });

  void _rows(NoteFolder folder, int depth, List<Widget> rows) {
    for (final child in folder.folders) {
      final open = !_collapsed.contains(child.path);
      rows.add(
        NoteTreeRow(
          key: ValueKey('folder:${child.path}'),
          depth: depth,
          icon: open ? Icons.folder_open_outlined : Icons.folder_outlined,
          label: child.name,
          strong: true,
          onTap: () => _toggle(child.path),
        ),
      );
      if (open) {
        _rows(child, depth + 1, rows);
      }
    }
    for (final note in folder.notes) {
      rows.add(
        NoteTreeRow(
          key: ValueKey('note:${note.id}'),
          depth: depth,
          icon: note.error != null
              ? Icons.sync_problem
              : note.meta.isCanvas
              ? Icons.dashboard_outlined
              : Icons.description_outlined,
          label: note.meta.label,
          selected: note.id == widget.active,
          onTap: () => widget.onOpen(note),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final rows = <Widget>[];
    _rows(widget.root, 0, rows);
    return ListView(children: rows);
  }
}
