import 'package:flutter/services.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';

class NoteTreeRow extends StatelessWidget {
  const NoteTreeRow({
    required this.depth,
    required this.icon,
    required this.label,
    required this.onTap,
    this.selected = false,
    this.strong = false,
    this.hovered = false,
    this.editing = false,
    this.onRename,
    this.onMenu,
    super.key,
  });

  static const indent = 16.0;

  final int depth;
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool selected;
  final bool strong;
  final bool hovered;
  final bool editing;
  final ValueChanged<String?>? onRename;
  final ValueChanged<Offset>? onMenu;

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    final style = Theme.of(context).textTheme.bodyMedium!.copyWith(
      color: selected ? palette.accent : palette.text,
      fontWeight: strong || selected ? FontWeight.w600 : null,
    );
    final menu = onMenu;
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: Spacing.s),
      child: Material(
        color: selected
            ? palette.accentSoft
            : hovered
            ? palette.accent.withValues(alpha: 0.12)
            : Colors.transparent,
        borderRadius: context.tokens.borderRadiusSmall,
        child: InkWell(
          borderRadius: context.tokens.borderRadiusSmall,
          onTap: editing ? null : onTap,
          onSecondaryTapDown: menu == null
              ? null
              : (details) => menu(details.globalPosition),
          onLongPress: menu == null
              ? null
              : () {
                  final box = context.findRenderObject()! as RenderBox;
                  menu(box.localToGlobal(box.size.centerLeft(Offset.zero)));
                },
          child: SizedBox(
            height: 30,
            child: Row(
              children: [
                for (var level = 0; level < depth; level++)
                  Container(
                    width: indent,
                    margin: const EdgeInsets.only(left: Spacing.s - 1),
                    decoration: BoxDecoration(
                      border: Border(left: BorderSide(color: palette.border)),
                    ),
                  ),
                const SizedBox(width: Spacing.s),
                Icon(
                  icon,
                  size: 16,
                  color: selected ? palette.accent : palette.muted,
                ),
                const SizedBox(width: Spacing.s),
                Expanded(
                  child: editing
                      ? _NameField(
                          label: label,
                          style: style,
                          onDone: onRename!,
                        )
                      : Text(
                          label,
                          style: style,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _NameField extends StatefulWidget {
  const _NameField({
    required this.label,
    required this.style,
    required this.onDone,
  });

  final String label;
  final TextStyle style;
  final ValueChanged<String?> onDone;

  @override
  State<_NameField> createState() => _NameFieldState();
}

class _NameFieldState extends State<_NameField> {
  late final _controller = TextEditingController(text: widget.label)
    ..selection = TextSelection(
      baseOffset: 0,
      extentOffset: widget.label.length,
    );
  final _focus = FocusNode();
  bool _done = false;

  @override
  void initState() {
    super.initState();
    _focus.addListener(() {
      if (!_focus.hasFocus) {
        _finish(_controller.text);
      }
    });
    WidgetsBinding.instance.addPostFrameCallback((_) => _focus.requestFocus());
  }

  void _finish(String? name) {
    if (_done) {
      return;
    }
    _done = true;
    widget.onDone(name);
  }

  @override
  void dispose() {
    _focus.dispose();
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return CallbackShortcuts(
      bindings: {
        const SingleActivator(LogicalKeyboardKey.escape): () => _finish(null),
      },
      child: TextField(
        controller: _controller,
        focusNode: _focus,
        style: widget.style,
        onSubmitted: _finish,
        decoration: const InputDecoration(
          isDense: true,
          contentPadding: EdgeInsets.symmetric(horizontal: 4, vertical: 4),
        ),
      ),
    );
  }
}
