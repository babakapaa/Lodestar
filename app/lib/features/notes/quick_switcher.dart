import 'package:flutter/services.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_resolver.dart';
import 'quick_switch_choice.dart';

class QuickSwitcher extends StatefulWidget {
  const QuickSwitcher({required this.notes, super.key});

  final List<Note> notes;

  static Future<QuickSwitchChoice?> show(
    BuildContext context,
    List<Note> notes,
  ) => showDialog<QuickSwitchChoice>(
    context: context,
    builder: (context) => QuickSwitcher(notes: notes),
  );

  @override
  State<QuickSwitcher> createState() => _QuickSwitcherState();
}

class _QuickSwitcherState extends State<QuickSwitcher> {
  String _query = '';
  int _selected = 0;

  List<Note> get _matches =>
      NoteResolver.search(_query, widget.notes).take(8).toList();

  bool get _canCreate =>
      _query.trim().isNotEmpty &&
      !widget.notes.any(
        (note) => !note.deleted && note.meta.name == _query.trim(),
      );

  int get _count => _matches.length + (_canCreate ? 1 : 0);

  void _choose(int index) {
    final matches = _matches;
    Navigator.of(context).pop(
      index < matches.length
          ? QuickSwitchChoice.open(matches[index].id)
          : QuickSwitchChoice.create(_query.trim()),
    );
  }

  KeyEventResult _onKey(FocusNode node, KeyEvent event) {
    if (event is! KeyDownEvent && event is! KeyRepeatEvent || _count == 0) {
      return KeyEventResult.ignored;
    }
    final key = event.logicalKey;
    if (key == LogicalKeyboardKey.arrowDown) {
      setState(() => _selected = (_selected + 1) % _count);
    } else if (key == LogicalKeyboardKey.arrowUp) {
      setState(() => _selected = (_selected - 1 + _count) % _count);
    } else {
      return KeyEventResult.ignored;
    }
    return KeyEventResult.handled;
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = context.palette;
    final matches = _matches;
    final selected = _selected.clamp(0, _count == 0 ? 0 : _count - 1);
    return Dialog(
      alignment: const Alignment(0, -0.6),
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 560, maxHeight: 460),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Padding(
              padding: const EdgeInsets.all(Spacing.m),
              child: Focus(
                onKeyEvent: _onKey,
                child: TextField(
                  autofocus: true,
                  decoration: InputDecoration(
                    prefixIcon: const Icon(Icons.search),
                    hintText: l10n.notesQuickHint,
                  ),
                  onChanged: (value) => setState(() {
                    _query = value;
                    _selected = 0;
                  }),
                  onSubmitted: (_) => _count == 0 ? null : _choose(selected),
                ),
              ),
            ),
            const Divider(),
            Flexible(
              child: ListView(
                shrinkWrap: true,
                padding: const EdgeInsets.symmetric(vertical: Spacing.s),
                children: [
                  for (final (index, note) in matches.indexed)
                    ListTile(
                      dense: true,
                      selected: index == selected,
                      selectedTileColor: palette.accentSoft,
                      leading: const Icon(Icons.description_outlined),
                      title: Text(note.meta.name),
                      subtitle: note.meta.folder.isEmpty
                          ? null
                          : Text(note.meta.folder),
                      onTap: () => _choose(index),
                    ),
                  if (_canCreate)
                    ListTile(
                      dense: true,
                      selected: selected == matches.length,
                      selectedTileColor: palette.accentSoft,
                      leading: const Icon(Icons.add),
                      title: Text(l10n.notesLinkCreate(_query.trim())),
                      onTap: () => _choose(matches.length),
                    ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
