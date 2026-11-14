import 'package:material_ui/material_ui.dart';

import '../../../design/tokens/palette.dart';
import 'canvas_board.dart';

class CanvasPainter extends CustomPainter {
  const CanvasPainter({
    required this.board,
    required this.pan,
    required this.scale,
    required this.palette,
    required this.selected,
    required this.dragLine,
  });

  static const grid = 24.0;

  final CanvasBoard board;
  final Offset pan;
  final double scale;
  final Palette palette;
  final String? selected;
  final (Offset, Offset)? dragLine;

  Offset _screen(Offset point) => pan + point * scale;

  static Path curve(Offset from, Offset to) {
    final bend = (to.dx - from.dx).abs() / 2 + 40;
    return Path()
      ..moveTo(from.dx, from.dy)
      ..cubicTo(from.dx + bend, from.dy, to.dx - bend, to.dy, to.dx, to.dy);
  }

  (Offset, Offset)? ends(String from, String to) {
    final a = board.card(from);
    final b = board.card(to);
    if (a == null || b == null) {
      return null;
    }
    return (_screen(a.rect.centerRight), _screen(b.rect.centerLeft));
  }

  @override
  void paint(Canvas canvas, Size size) {
    final dot = Paint()..color = palette.border;
    final step = grid * scale;
    if (step >= 8) {
      final startX = pan.dx % step;
      final startY = pan.dy % step;
      for (var x = startX; x < size.width; x += step) {
        for (var y = startY; y < size.height; y += step) {
          canvas.drawCircle(Offset(x, y), 1, dot);
        }
      }
    }
    for (final edge in board.edges) {
      final points = ends(edge.from, edge.to);
      if (points == null) {
        continue;
      }
      final active = edge.id == selected;
      _arrow(
        canvas,
        points.$1,
        points.$2,
        active ? palette.accent : palette.muted,
        active ? 3 : 2,
      );
    }
    final line = dragLine;
    if (line != null) {
      _arrow(canvas, line.$1, line.$2, palette.accent, 2);
    }
  }

  void _arrow(
    Canvas canvas,
    Offset from,
    Offset to,
    Color color,
    double width,
  ) {
    final stroke = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = width;
    canvas.drawPath(curve(from, to), stroke);
    final head = Path()
      ..moveTo(to.dx, to.dy)
      ..lineTo(to.dx - 10, to.dy - 6)
      ..lineTo(to.dx - 10, to.dy + 6)
      ..close();
    canvas.drawPath(head, Paint()..color = color);
  }

  String? edgeAt(Offset point) {
    for (final edge in board.edges.reversed) {
      final points = ends(edge.from, edge.to);
      if (points == null) {
        continue;
      }
      for (final metric in curve(points.$1, points.$2).computeMetrics()) {
        for (var at = 0.0; at <= metric.length; at += 6) {
          final tangent = metric.getTangentForOffset(at);
          if (tangent != null && (tangent.position - point).distance < 8) {
            return edge.id;
          }
        }
      }
    }
    return null;
  }

  @override
  bool shouldRepaint(CanvasPainter oldDelegate) => true;
}
