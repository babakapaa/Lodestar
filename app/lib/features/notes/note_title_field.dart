import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_meta.dart';
import 'notes_providers.dart';

class NoteTitleField extends ConsumerStatefulWidget {
  const NoteTitleField({
    required this.note,
    required this.readOnly,
    this.onDone,
    super.key,
  });

  final Note note;
  final bool readOnly;
  final VoidCallback? onDone;

  @override
  ConsumerState<NoteTitleField> createState() => _NoteTitleFieldState();
}

class _NoteTitleFieldState extends ConsumerState<NoteTitleField> {
  late final _controller = TextEditingController(text: widget.note.meta.label);
  final _focus = FocusNode();

  @override
  void initState() {
    super.initState();
    _focus.addListener(_onFocus);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted || widget.readOnly) {
        return;
      }
      if (widget.note.meta.label.startsWith(context.l10n.notesUntitled)) {
        _focus.requestFocus();
        _controller.selection = TextSelection(
          baseOffset: 0,
          extentOffset: _controller.text.length,
        );
      }
    });
  }

  @override
  void didUpdateWidget(NoteTitleField oldWidget) {
    super.didUpdateWidget(oldWidget);
    final label = widget.note.meta.label;
    if (!_focus.hasFocus && _controller.text != label) {
      _controller.text = label;
    }
  }

  void _onFocus() {
    if (!_focus.hasFocus) {
      _save();
    }
  }

  Future<void> _save() async {
    final meta = widget.note.meta;
    final name = _controller.text.replaceAll('/', ' ').trim();
    if (name.isEmpty) {
      _controller.text = meta.label;
      return;
    }
    final file = meta.isCanvas ? '$name${NoteMeta.canvasExtension}' : name;
    final path = NoteMeta.normalizePath(
      meta.folder.isEmpty ? file : '${meta.folder}/$file',
    );
    if (path != meta.path) {
      await ref
          .read(notesStoreProvider)
          .setMeta(widget.note.id, meta.copyWith(path: path));
    }
  }

  @override
  void dispose() {
    _focus
      ..removeListener(_onFocus)
      ..dispose();
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return TextField(
      key: const ValueKey('note-title'),
      controller: _controller,
      focusNode: _focus,
      readOnly: widget.readOnly,
      maxLines: null,
      textInputAction: TextInputAction.done,
      style: Theme.of(context).textTheme.headlineMedium,
      onSubmitted: (_) {
        _save();
        widget.onDone?.call();
      },
      decoration: InputDecoration(
        hintText: context.l10n.notesUntitled,
        border: InputBorder.none,
        enabledBorder: InputBorder.none,
        focusedBorder: InputBorder.none,
        filled: false,
        isCollapsed: true,
        contentPadding: const EdgeInsets.fromLTRB(
          Spacing.xl,
          Spacing.xl,
          Spacing.xl,
          Spacing.xs,
        ),
      ),
    );
  }
}
