import 'dart:math';

import 'package:material_ui/material_ui.dart';

import '../../design/tokens/palette.dart';
import 'notes_graph.dart';

class NotesGraphPainter extends CustomPainter {
  const NotesGraphPainter({
    required this.graph,
    required this.palette,
    required this.labelStyle,
    required this.origin,
  });

  static double radiusOf(int degree) => 6 + min(degree, 10) * 1.5;

  final NotesGraph graph;
  final Palette palette;
  final TextStyle labelStyle;
  final Offset origin;

  @override
  void paint(Canvas canvas, Size size) {
    canvas.translate(origin.dx, origin.dy);
    final edge = Paint()
      ..color = palette.border
      ..strokeWidth = 1.2;
    for (final (a, b) in graph.edges) {
      canvas.drawLine(graph.nodes[a].position, graph.nodes[b].position, edge);
    }
    final dot = Paint()..color = palette.accent;
    final lonely = Paint()..color = palette.muted;
    for (final node in graph.nodes) {
      final radius = radiusOf(node.degree);
      canvas.drawCircle(node.position, radius, node.degree > 0 ? dot : lonely);
      final label = TextPainter(
        text: TextSpan(text: node.label, style: labelStyle),
        textDirection: TextDirection.ltr,
        maxLines: 1,
        ellipsis: '…',
      )..layout(maxWidth: 160);
      label.paint(canvas, node.position + Offset(-label.width / 2, radius + 4));
    }
  }

  @override
  bool shouldRepaint(NotesGraphPainter oldDelegate) =>
      oldDelegate.graph != graph ||
      oldDelegate.palette != palette ||
      oldDelegate.origin != origin;
}
