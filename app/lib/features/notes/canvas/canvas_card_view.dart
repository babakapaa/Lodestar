import 'package:material_ui/material_ui.dart';

import '../../../design/tokens/spacing.dart';
import '../../../design/tokens/tokens_context.dart';
import '../markdown_controller.dart';
import 'canvas_card.dart';

class CanvasCardView extends StatefulWidget {
  const CanvasCardView({
    required this.card,
    required this.title,
    required this.preview,
    required this.selected,
    required this.editing,
    required this.readOnly,
    required this.onSelect,
    required this.onActivate,
    required this.onMove,
    required this.onMoved,
    required this.onResize,
    required this.onText,
    required this.onConnectStart,
    required this.onConnectUpdate,
    required this.onConnectEnd,
    super.key,
  });

  final CanvasCard card;
  final String? title;
  final String preview;
  final bool selected;
  final bool editing;
  final bool readOnly;
  final VoidCallback onSelect;
  final VoidCallback onActivate;
  final ValueChanged<Offset> onMove;
  final VoidCallback onMoved;
  final ValueChanged<Offset> onResize;
  final ValueChanged<String> onText;
  final ValueChanged<Offset> onConnectStart;
  final ValueChanged<Offset> onConnectUpdate;
  final ValueChanged<Offset> onConnectEnd;

  @override
  State<CanvasCardView> createState() => _CanvasCardViewState();
}

class _CanvasCardViewState extends State<CanvasCardView> {
  late final _controller = MarkdownController()..text = widget.card.text;
  Offset _pointer = Offset.zero;

  @override
  void didUpdateWidget(CanvasCardView oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (!widget.editing && _controller.text != widget.card.text) {
      _controller.text = widget.card.text;
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    final textTheme = Theme.of(context).textTheme;
    final body = textTheme.bodyMedium!.copyWith(height: 1.5);
    final title = widget.title;
    final content = widget.editing
        ? TextField(
            controller: _controller,
            autofocus: true,
            maxLines: null,
            expands: true,
            style: body,
            onChanged: widget.onText,
            decoration: const InputDecoration(
              border: InputBorder.none,
              enabledBorder: InputBorder.none,
              focusedBorder: InputBorder.none,
              filled: false,
              isCollapsed: true,
            ),
          )
        : ClipRect(
            child: Text.rich(
              _controller.buildTextSpan(
                context: context,
                style: body,
                withComposing: false,
              ),
              overflow: TextOverflow.fade,
            ),
          );
    return Stack(
      clipBehavior: Clip.none,
      children: [
        Positioned.fill(
          child: GestureDetector(
            onTap: widget.onSelect,
            onDoubleTap: widget.onActivate,
            onPanStart: widget.readOnly || widget.editing
                ? null
                : (_) => widget.onSelect(),
            onPanUpdate: widget.readOnly || widget.editing
                ? null
                : (details) => widget.onMove(details.delta),
            onPanEnd: widget.readOnly || widget.editing
                ? null
                : (_) => widget.onMoved(),
            child: Container(
              decoration: BoxDecoration(
                color: palette.surface,
                borderRadius: tokens.borderRadiusSmall,
                border: Border.all(
                  color: widget.selected ? palette.accent : palette.border,
                  width: widget.selected ? 2 : 1,
                ),
                boxShadow: palette.shadow,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  if (title != null)
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: Spacing.m,
                        vertical: Spacing.s,
                      ),
                      decoration: BoxDecoration(
                        color: palette.accentSoft,
                        border: Border(
                          bottom: BorderSide(color: palette.border),
                        ),
                      ),
                      child: Row(
                        children: [
                          Icon(
                            Icons.description_outlined,
                            size: 16,
                            color: palette.accent,
                          ),
                          const SizedBox(width: Spacing.s),
                          Expanded(
                            child: Text(
                              title,
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: textTheme.titleSmall!.copyWith(
                                color: palette.accent,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  Expanded(
                    child: Padding(
                      padding: const EdgeInsets.all(Spacing.m),
                      child: title != null
                          ? Text(
                              widget.preview,
                              overflow: TextOverflow.fade,
                              style: body.copyWith(color: palette.muted),
                            )
                          : content,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        if (!widget.readOnly && widget.selected) ...[
          Positioned(
            right: -7,
            top: 0,
            bottom: 0,
            child: Center(
              child: GestureDetector(
                onPanStart: (details) {
                  _pointer = details.globalPosition;
                  widget.onConnectStart(_pointer);
                },
                onPanUpdate: (details) {
                  _pointer = details.globalPosition;
                  widget.onConnectUpdate(_pointer);
                },
                onPanEnd: (_) => widget.onConnectEnd(_pointer),
                child: Container(
                  width: 14,
                  height: 14,
                  decoration: BoxDecoration(
                    color: palette.accent,
                    shape: BoxShape.circle,
                    border: Border.all(color: palette.surface, width: 2),
                  ),
                ),
              ),
            ),
          ),
          Positioned(
            right: 0,
            bottom: 0,
            child: GestureDetector(
              onPanUpdate: (details) => widget.onResize(details.delta),
              onPanEnd: (_) => widget.onMoved(),
              child: Icon(Icons.south_east, size: 16, color: palette.accent),
            ),
          ),
        ],
      ],
    );
  }
}
