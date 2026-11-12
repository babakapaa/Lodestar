import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'note_doc.dart';
import 'notes_providers.dart';
import 'notes_store.dart';
import 'text_edit.dart';

class NoteEditor extends ConsumerStatefulWidget {
  const NoteEditor({required this.noteId, required this.readOnly, super.key});

  final String noteId;
  final bool readOnly;

  @override
  ConsumerState<NoteEditor> createState() => _NoteEditorState();
}

class _NoteEditorState extends ConsumerState<NoteEditor> {
  final _controller = TextEditingController();
  late final NotesStore _store;
  NoteDoc? _doc;
  StreamSubscription<TextEdit>? _remote;

  @override
  void initState() {
    super.initState();
    _store = ref.read(notesStoreProvider);
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
    _remote = doc.remoteEdits.listen(_onRemote);
    _controller.addListener(_onChanged);
  }

  void _onChanged() {
    final doc = _doc!;
    if (_controller.text != doc.text) {
      _store.edit(doc, _controller.text);
    }
  }

  void _onRemote(TextEdit change) {
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

  @override
  void dispose() {
    _controller
      ..removeListener(_onChanged)
      ..dispose();
    unawaited(_remote?.cancel());
    unawaited(_store.flush(widget.noteId));
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (_doc == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return TextField(
      controller: _controller,
      readOnly: widget.readOnly,
      maxLines: null,
      expands: true,
      keyboardType: TextInputType.multiline,
      textAlignVertical: TextAlignVertical.top,
      decoration: InputDecoration(
        hintText: context.l10n.notesEditorHint,
        border: InputBorder.none,
        enabledBorder: InputBorder.none,
        focusedBorder: InputBorder.none,
        filled: false,
        contentPadding: const EdgeInsets.all(Spacing.l),
      ),
    );
  }
}
