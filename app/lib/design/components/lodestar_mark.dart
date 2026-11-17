import 'package:material_ui/material_ui.dart';

import 'lodestar_mark_painter.dart';

class LodestarMark extends StatelessWidget {
  const LodestarMark({required this.size, required this.color, super.key});

  final double size;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return SizedBox.square(
      dimension: size,
      child: CustomPaint(painter: LodestarMarkPainter(color: color)),
    );
  }
}
