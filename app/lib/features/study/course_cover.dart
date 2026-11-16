import 'package:material_ui/material_ui.dart';

abstract final class CourseCover {
  static List<Color> of(String seed, bool dark) {
    var hash = 0;
    for (final unit in seed.codeUnits) {
      hash = (hash * 31 + unit) & 0x7fffffff;
    }
    final hue = (hash % 360).toDouble();
    final lightness = dark ? 0.32 : 0.55;
    return [
      HSLColor.fromAHSL(1, hue, 0.55, lightness).toColor(),
      HSLColor.fromAHSL(1, (hue + 40) % 360, 0.6, lightness + 0.08).toColor(),
    ];
  }
}
