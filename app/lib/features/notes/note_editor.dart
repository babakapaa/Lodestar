import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'link_suggestions.dart';
import 'markdown_controller.dart';
import 'note.dart';
import 'note_doc.dart';
import 'note_resolver.dart';
import 'notes_providers.dart';
import 'notes_store.dart';
import 'text_edit.dart';
import 'wiki_link_query.dart';

class NoteEditor extends ConsumerStatefulWidget {
  const NoteEditor({
    required this.noteId,
    required this.readOnly,
    this.onOpenLink,
    this.jumps,
    super.key,
  });

  static const readableWidth = 760.0;

  final String noteId;
  final bool readOnly;
  final ValueChanged<String>? onOpenLink;
  final ValueListenable<int?>? jumps;

  @override
  ConsumerState<NoteEditor> createState() => _NoteEditorState();
}

class _NoteEditorState extends ConsumerState<NoteEditor> {
  final _controller = MarkdownController();
  late final _focus = FocusNode(onKeyEvent: _onKey);
  late final NotesStore _store;
  NoteDoc? _doc;
  StreamSubscription<TextEdit>? _outside;
  String? _query;
  int _selected = 0;
  bool _dismissed = false;

  @override
  void initState() {
    super.initState();
    _store = ref.read(notesStoreProvider);
    widget.jumps?.addListener(_jump);
    unawaited(_open());
  }

  Future<void> _open() async {
    final doc = await _store.open(widget.noteId);
    if (!mounted) {
      return;
    }
    setState(() {
      _doc = doc;
      _controller.text = doc.text;
    });
    _outside = doc.outsideEdits.listen(_onOutside);
    _controller.addListener(_onChanged);
  }

  void _onChanged() {
    final doc = _doc!;
    if (_controller.text != doc.text) {
      _store.edit(doc, _controller.text);
    }
    final query = widget.readOnly
        ? null
        : WikiLinkQuery.typed(_controller.value);
    if (query != _query) {
      setState(() {
        _query = query;
        _selected = 0;
        _dismissed = false;
      });
    }
  }

  void _onOutside(TextEdit change) {
    final value = _controller.value;
    final selection = value.selection;
    _controller.value = TextEditingValue(
      text: _doc!.text,
      selection: selection.isValid
          ? TextSelection(
              baseOffset: change.shift(selection.baseOffset),
              extentOffset: change.shift(selection.extentOffset),
            )
          : selection,
    );
  }

  void _jump() {
    final offset = widget.jumps?.value;
    if (offset == null || _doc == null) {
      return;
    }
    _controller.selection = TextSelection.collapsed(
      offset: offset.clamp(0, _controller.text.length),
    );
    _focus.requestFocus();
  }

  List<Note> _matches() {
    final query = _query;
    if (query == null || _dismissed) {
      return const [];
    }
    final notes = ref.read(notesProvider).value ?? const <Note>[];
    return NoteResolver.search(
      query,
      notes.where((note) => note.id != widget.noteId),
    ).take(6).toList();
  }

  void _pick(Note note) {
    _controller.value = WikiLinkQuery.insert(_controller.value, note.meta.name);
    _focus.requestFocus();
  }

  KeyEventResult _onKey(FocusNode node, KeyEvent event) {
    if (event is! KeyDownEvent && event is! KeyRepeatEvent) {
      return KeyEventResult.ignored;
    }
    final key = event.logicalKey;
    final keyboard = HardwareKeyboard.instance;
    if (key == LogicalKeyboardKey.enter &&
        (keyboard.isControlPressed || keyboard.isMetaPressed)) {
      final link = WikiLinkQuery.at(_controller.value);
      if (link != null && widget.onOpenLink != null) {
        widget.onOpenLink!(link);
        return KeyEventResult.handled;
      }
      return KeyEventResult.ignored;
    }
    final matches = _matches();
    if (matches.isEmpty) {
      return KeyEventResult.ignored;
    }
    if (key == LogicalKeyboardKey.arrowDown) {
      setState(() => _selected = (_selected + 1) % matches.length);
    } else if (key == LogicalKeyboardKey.arrowUp) {
      setState(
        () => _selected = (_selected - 1 + matches.length) % matches.length,
      );
    } else if (key == LogicalKeyboardKey.enter ||
        key == LogicalKeyboardKey.tab) {
      _pick(matches[_selected.clamp(0, matches.length - 1)]);
    } else if (key == LogicalKeyboardKey.escape) {
      setState(() => _dismissed = true);
    } else {
      return KeyEventResult.ignored;
    }
    return KeyEventResult.handled;
  }

  @override
  void didUpdateWidget(NoteEditor oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.jumps != widget.jumps) {
      oldWidget.jumps?.removeListener(_jump);
      widget.jumps?.addListener(_jump);
    }
  }

  @override
  void dispose() {
    widget.jumps?.removeListener(_jump);
    _controller
      ..removeListener(_onChanged)
      ..dispose();
    _focus.dispose();
    unawaited(_outside?.cancel());
    unawaited(_store.flush(widget.noteId));
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (_doc == null) {
      return const Center(child: CircularProgressIndicator());
    }
    final matches = _matches();
    final body = Theme.of(context).textTheme.bodyLarge!;
    return Stack(
      children: [
        Positioned.fill(
          child: Align(
            alignment: Alignment.topCenter,
            child: ConstrainedBox(
              constraints: const BoxConstraints(
                maxWidth: NoteEditor.readableWidth,
              ),
              child: TextField(
                controller: _controller,
                focusNode: _focus,
                readOnly: widget.readOnly,
                maxLines: null,
                expands: true,
                keyboardType: TextInputType.multiline,
                textAlignVertical: TextAlignVertical.top,
                style: body.copyWith(height: 1.6),
                decoration: InputDecoration(
                  hintText: context.l10n.notesEditorHint,
                  border: InputBorder.none,
                  enabledBorder: InputBorder.none,
                  focusedBorder: InputBorder.none,
                  filled: false,
                  contentPadding: const EdgeInsets.fromLTRB(
                    Spacing.xl,
                    Spacing.l,
                    Spacing.xl,
                    Spacing.xl,
                  ),
                ),
              ),
            ),
          ),
        ),
        if (matches.isNotEmpty)
          Positioned(
            left: Spacing.l,
            right: Spacing.l,
            bottom: Spacing.l,
            child: Align(
              alignment: Alignment.bottomCenter,
              child: ConstrainedBox(
                constraints: const BoxConstraints(maxWidth: 420),
                child: LinkSuggestions(
                  notes: matches,
                  selected: _selected.clamp(0, matches.length - 1),
                  onPick: _pick,
                ),
              ),
            ),
          ),
      ],
    );
  }
}
