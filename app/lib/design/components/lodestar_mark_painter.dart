import 'package:material_ui/material_ui.dart';

class LodestarMarkPainter extends CustomPainter {
  const LodestarMarkPainter({required this.color});

  final Color color;

  static const _grid = 512.0;
  static const _radius = 116.0;

  @override
  void paint(Canvas canvas, Size size) {
    final scale = size.shortestSide / _grid;
    canvas.scale(scale);
    final tile = Path()
      ..addRRect(
        RRect.fromRectAndRadius(
          const Rect.fromLTWH(0, 0, _grid, _grid),
          const Radius.circular(_radius),
        ),
      );
    final star = Path()
      ..moveTo(256, 0)
      ..cubicTo(264, 216, 296, 248, 512, 256)
      ..cubicTo(296, 264, 264, 296, 256, 512)
      ..cubicTo(248, 296, 216, 264, 0, 256)
      ..cubicTo(216, 248, 248, 216, 256, 0)
      ..close();
    canvas.drawPath(
      Path.combine(PathOperation.difference, tile, star),
      Paint()
        ..color = color
        ..isAntiAlias = true,
    );
  }

  @override
  bool shouldRepaint(LodestarMarkPainter old) => old.color != color;
}
