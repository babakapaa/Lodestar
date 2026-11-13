import 'dart:math';

import 'package:material_ui/material_ui.dart';

import '../../design/tokens/palette.dart';
import 'search_map.dart';

class SearchMapPainter extends CustomPainter {
  const SearchMapPainter({
    required this.map,
    required this.center,
    required this.palette,
    required this.labelStyle,
  });

  static const _labelWidth = 120.0;

  final SearchMap map;
  final String center;
  final Palette palette;
  final TextStyle labelStyle;

  @override
  void paint(Canvas canvas, Size size) {
    final middle = size.center(Offset.zero);
    final reach = min(size.width, size.height) / 2 - 24;
    final used = _ring(map.used.length, middle, reach * 0.55, 0);
    final rejected = _ring(map.rejected.length, middle, reach, pi / 7);
    final solid = Paint()
      ..color = palette.accent
      ..strokeWidth = 1.6;
    final faint = Paint()
      ..color = palette.border
      ..strokeWidth = 1;
    for (final point in used) {
      canvas.drawLine(middle, point, solid);
    }
    for (final point in rejected) {
      canvas.drawLine(middle, point, faint);
    }
    canvas.drawCircle(middle, 14, Paint()..color = palette.text);
    _label(canvas, center, middle, 14, labelStyle);
    for (final (index, point) in used.indexed) {
      canvas.drawCircle(
        point,
        8,
        Paint()
          ..color = map.used[index].scholarly
              ? palette.accent
              : palette.accentSoft,
      );
      _label(canvas, _short(map.used[index].reference), point, 8, labelStyle);
    }
    final muted = labelStyle.copyWith(color: palette.muted);
    for (final (index, point) in rejected.indexed) {
      canvas.drawCircle(point, 6, Paint()..color = palette.muted);
      _label(canvas, _short(map.rejected[index].reference), point, 6, muted);
    }
  }

  static List<Offset> _ring(
    int count,
    Offset middle,
    double radius,
    double turn,
  ) => [
    for (var i = 0; i < count; i++)
      middle + Offset.fromDirection(turn - pi / 2 + 2 * pi * i / count, radius),
  ];

  static String _short(String reference) {
    final cut = reference.indexOf(RegExp(r'[./]'), 1);
    return cut > 0 ? reference.substring(0, cut) : reference;
  }

  static void _label(
    Canvas canvas,
    String text,
    Offset point,
    double radius,
    TextStyle style,
  ) {
    final label = TextPainter(
      text: TextSpan(text: text, style: style),
      textDirection: TextDirection.ltr,
      textAlign: TextAlign.center,
      maxLines: 2,
      ellipsis: '…',
    )..layout(maxWidth: _labelWidth);
    label.paint(canvas, point + Offset(-label.width / 2, radius + 4));
  }

  @override
  bool shouldRepaint(SearchMapPainter oldDelegate) =>
      oldDelegate.map != map ||
      oldDelegate.center != center ||
      oldDelegate.palette != palette;
}
