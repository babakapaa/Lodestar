import 'dart:async';
import 'dart:math';

import 'package:flutter/gestures.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../../design/tokens/spacing.dart';
import '../../../design/tokens/tokens_context.dart';
import '../../../l10n/l10n_context.dart';
import '../note.dart';
import '../note_doc.dart';
import '../note_resolver.dart';
import '../notes_providers.dart';
import '../notes_store.dart';
import '../quick_switcher.dart';
import '../text_edit.dart';
import 'canvas_board.dart';
import 'canvas_card.dart';
import 'canvas_card_view.dart';
import 'canvas_painter.dart';
import 'canvas_toolbar.dart';

class CanvasView extends ConsumerStatefulWidget {
  const CanvasView({
    required this.noteId,
    required this.readOnly,
    required this.onOpenNote,
    super.key,
  });

  final String noteId;
  final bool readOnly;
  final ValueChanged<String> onOpenNote;

  @override
  ConsumerState<CanvasView> createState() => _CanvasViewState();
}

class _CanvasViewState extends ConsumerState<CanvasView> {
  final _stack = GlobalKey();
  final _focus = FocusNode();
  late final NotesStore _store;
  NoteDoc? _doc;
  StreamSubscription<TextEdit>? _outside;
  CanvasBoard _board = CanvasBoard([], []);
  Offset _pan = const Offset(Spacing.xl, Spacing.xl);
  double _scale = 1;
  Size _viewport = Size.zero;
  String? _selected;
  String? _editing;
  (String, Offset)? _connecting;
  bool _fitted = false;

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
      _board = CanvasBoard.parse(doc.text);
    });
    _outside = doc.outsideEdits.listen(
      (_) => setState(() => _board = CanvasBoard.parse(doc.text)),
    );
  }

  void _save() {
    final doc = _doc;
    if (doc == null || widget.readOnly) {
      return;
    }
    final text = _board.encode();
    if (text != doc.text) {
      _store.edit(doc, text);
    }
  }

  Offset _toBoard(Offset local) => (local - _pan) / _scale;

  Offset _local(Offset global) {
    final box = _stack.currentContext!.findRenderObject()! as RenderBox;
    return box.globalToLocal(global);
  }

  void _addCard(Offset at, {String? file}) {
    final card = CanvasCard(
      id: CanvasBoard.newId(),
      rect: Rect.fromLTWH(at.dx, at.dy, 260, file == null ? 120 : 180),
      file: file,
    );
    setState(() {
      _board.cards.add(card);
      _selected = card.id;
      _editing = file == null ? card.id : null;
    });
    _save();
  }

  Offset get _center =>
      _toBoard(Offset(_viewport.width / 2 - 130, _viewport.height / 2 - 60));

  Future<void> _addNote() async {
    final notes = ref.read(notesProvider).value ?? const <Note>[];
    final choice = await QuickSwitcher.show(
      context,
      notes.where((note) => !note.meta.isCanvas).toList(),
    );
    if (choice == null || !mounted) {
      return;
    }
    final id = choice.noteId;
    final path = id == null
        ? choice.createPath!
        : notes.firstWhere((note) => note.id == id).meta.path;
    _addCard(_center, file: '$path.md');
  }

  void _zoom(double factor, [Offset? around]) {
    final focus = around ?? _viewport.center(Offset.zero);
    final next = (_scale * factor).clamp(0.2, 3.0);
    setState(() {
      _pan = focus - (focus - _pan) * (next / _scale);
      _scale = next;
    });
  }

  void _fit() {
    final bounds = _board.bounds;
    if (bounds == null || _viewport.isEmpty) {
      return;
    }
    final area = bounds.inflate(48);
    final scale = min(
      1.0,
      min(_viewport.width / area.width, _viewport.height / area.height),
    ).clamp(0.2, 1.0);
    setState(() {
      _scale = scale;
      _pan = _viewport.center(Offset.zero) - area.center * scale;
    });
  }

  void _deleteSelected() {
    final id = _selected;
    if (id == null || widget.readOnly) {
      return;
    }
    setState(() {
      _board.remove(id);
      _selected = null;
    });
    _save();
  }

  KeyEventResult _onKey(FocusNode node, KeyEvent event) {
    if (event is! KeyDownEvent || _editing != null) {
      if (event is KeyDownEvent &&
          event.logicalKey == LogicalKeyboardKey.escape) {
        setState(() => _editing = null);
        return KeyEventResult.handled;
      }
      return KeyEventResult.ignored;
    }
    if (event.logicalKey == LogicalKeyboardKey.delete ||
        event.logicalKey == LogicalKeyboardKey.backspace) {
      _deleteSelected();
      return KeyEventResult.handled;
    }
    return KeyEventResult.ignored;
  }

  void _tapBackground(Offset local) {
    final painter = _painter(context);
    setState(() {
      _selected = painter.edgeAt(local);
      _editing = null;
    });
    _focus.requestFocus();
  }

  CanvasPainter _painter(BuildContext context) => CanvasPainter(
    board: _board,
    pan: _pan,
    scale: _scale,
    palette: context.palette,
    selected: _selected,
    dragLine: switch (_connecting) {
      (final from, final pointer) when _board.card(from) != null => (
        _pan + _board.card(from)!.rect.centerRight * _scale,
        pointer,
      ),
      _ => null,
    },
  );

  @override
  void dispose() {
    unawaited(_outside?.cancel());
    _focus.dispose();
    unawaited(_store.flush(widget.noteId));
    super.dispose();
  }

  Widget _card(CanvasCard card, List<Note> notes) {
    final path = card.notePath;
    final target = path == null ? null : NoteResolver.resolve(path, notes);
    final rect = card.rect;
    return Positioned(
      key: ValueKey('card:${card.id}'),
      left: _pan.dx + rect.left * _scale,
      top: _pan.dy + rect.top * _scale,
      width: rect.width * _scale,
      height: rect.height * _scale,
      child: FittedBox(
        fit: BoxFit.fill,
        alignment: Alignment.topLeft,
        child: SizedBox(
          width: rect.width,
          height: rect.height,
          child: CanvasCardView(
            card: card,
            title: path == null
                ? null
                : (target?.meta.name ?? path.split('/').last),
            preview: target?.body ?? '',
            selected: card.id == _selected,
            editing: card.id == _editing,
            readOnly: widget.readOnly,
            onSelect: () {
              setState(() => _selected = card.id);
              if (_editing != card.id) {
                _focus.requestFocus();
              }
            },
            onActivate: () {
              if (path != null) {
                widget.onOpenNote(path);
              } else if (!widget.readOnly) {
                setState(() {
                  _selected = card.id;
                  _editing = card.id;
                });
              }
            },
            onMove: (delta) => setState(() => card.rect = rect.shift(delta)),
            onMoved: _save,
            onResize: (delta) => setState(() {
              card.rect = Rect.fromLTWH(
                card.rect.left,
                card.rect.top,
                max(CanvasCard.minSize.width, card.rect.width + delta.dx),
                max(CanvasCard.minSize.height, card.rect.height + delta.dy),
              );
            }),
            onText: (text) {
              card.text = text;
              _save();
            },
            onConnectStart: (global) =>
                setState(() => _connecting = (card.id, _local(global))),
            onConnectUpdate: (global) =>
                setState(() => _connecting = (card.id, _local(global))),
            onConnectEnd: (global) {
              final target = _board.cardAt(_toBoard(_local(global)));
              setState(() {
                if (target != null) {
                  _board.connect(card.id, target.id);
                }
                _connecting = null;
              });
              _save();
            },
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (_doc == null) {
      return const Center(child: CircularProgressIndicator());
    }
    final notes = ref.watch(notesProvider).value ?? const <Note>[];
    final palette = context.palette;
    return LayoutBuilder(
      builder: (context, constraints) {
        _viewport = constraints.biggest;
        if (!_fitted) {
          _fitted = true;
          WidgetsBinding.instance.addPostFrameCallback((_) {
            if (mounted) {
              _fit();
            }
          });
        }
        return Focus(
          focusNode: _focus,
          onKeyEvent: _onKey,
          child: Listener(
            onPointerSignal: (event) {
              if (event is! PointerScrollEvent) {
                return;
              }
              final keyboard = HardwareKeyboard.instance;
              if (keyboard.isControlPressed || keyboard.isMetaPressed) {
                _zoom(
                  event.scrollDelta.dy > 0 ? 0.9 : 1.1,
                  event.localPosition,
                );
              } else {
                setState(() => _pan -= event.scrollDelta);
              }
            },
            child: ClipRect(
              child: Stack(
                key: _stack,
                children: [
                  Positioned.fill(
                    child: GestureDetector(
                      onTapUp: (details) =>
                          _tapBackground(details.localPosition),
                      onDoubleTapDown: widget.readOnly
                          ? null
                          : (details) =>
                                _addCard(_toBoard(details.localPosition)),
                      onDoubleTap: () {},
                      onPanUpdate: (details) =>
                          setState(() => _pan += details.delta),
                      child: ColoredBox(
                        color: palette.bg,
                        child: CustomPaint(painter: _painter(context)),
                      ),
                    ),
                  ),
                  for (final card in _board.cards) _card(card, notes),
                  if (_board.cards.isEmpty)
                    IgnorePointer(
                      child: Center(
                        child: Padding(
                          padding: const EdgeInsets.all(Spacing.xl),
                          child: Text(
                            context.l10n.canvasHint,
                            textAlign: TextAlign.center,
                            style: Theme.of(context).textTheme.bodyMedium!
                                .copyWith(color: palette.muted),
                          ),
                        ),
                      ),
                    ),
                  Positioned(
                    left: 0,
                    right: 0,
                    bottom: Spacing.l,
                    child: Center(
                      child: CanvasToolbar(
                        readOnly: widget.readOnly,
                        onAddCard: () => _addCard(_center),
                        onAddNote: _addNote,
                        onZoomIn: () => _zoom(1.2),
                        onZoomOut: () => _zoom(1 / 1.2),
                        onFit: _fit,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
