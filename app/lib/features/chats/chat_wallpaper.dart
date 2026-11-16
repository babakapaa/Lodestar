import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';

class ChatWallpaper extends StatelessWidget {
  const ChatWallpaper({required this.child, super.key});

  final Widget child;

  @override
  Widget build(BuildContext context) {
    final palette = context.tokens.palette;
    final base = Color.alphaBlend(
      palette.accent.withValues(alpha: palette.isDark ? 0.05 : 0.04),
      palette.side,
    );
    return CustomPaint(
      painter: _Dots(
        base: base,
        dot: palette.accent.withValues(alpha: palette.isDark ? 0.10 : 0.09),
      ),
      child: child,
    );
  }
}

class _Dots extends CustomPainter {
  const _Dots({required this.base, required this.dot});

  final Color base;
  final Color dot;

  static const step = 22.0;

  @override
  void paint(Canvas canvas, Size size) {
    canvas.drawRect(Offset.zero & size, Paint()..color = base);
    final paint = Paint()..color = dot;
    var row = 0;
    for (var y = step / 2; y < size.height; y += step) {
      final shift = row.isOdd ? step / 2 : 0.0;
      for (var x = step / 2 + shift; x < size.width; x += step) {
        canvas.drawCircle(Offset(x, y), 1.2, paint);
      }
      row++;
    }
  }

  @override
  bool shouldRepaint(_Dots old) => old.base != base || old.dot != dot;
}
