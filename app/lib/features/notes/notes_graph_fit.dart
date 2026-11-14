import 'dart:math';

import 'package:material_ui/material_ui.dart';

class NotesGraphFit extends StatefulWidget {
  const NotesGraphFit({
    required this.viewport,
    required this.content,
    required this.child,
    super.key,
  });

  final Size viewport;
  final Size content;
  final Widget child;

  @override
  State<NotesGraphFit> createState() => _NotesGraphFitState();
}

class _NotesGraphFitState extends State<NotesGraphFit> {
  final _transform = TransformationController();
  Size? _fitted;

  void _fit() {
    final viewport = widget.viewport;
    final content = widget.content;
    if (_fitted == viewport || content.isEmpty || viewport.isEmpty) {
      return;
    }
    _fitted = viewport;
    final scale = min(
      1.0,
      min(viewport.width / content.width, viewport.height / content.height),
    );
    final dx = (viewport.width - content.width * scale) / 2;
    final dy = (viewport.height - content.height * scale) / 2;
    _transform.value = Matrix4.identity()
      ..translateByDouble(dx, dy, 0, 1)
      ..scaleByDouble(scale, scale, 1, 1);
  }

  @override
  void dispose() {
    _transform.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    _fit();
    return InteractiveViewer(
      transformationController: _transform,
      constrained: false,
      boundaryMargin: const EdgeInsets.all(double.infinity),
      minScale: 0.2,
      maxScale: 3,
      child: widget.child,
    );
  }
}
