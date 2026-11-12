import 'package:material_ui/material_ui.dart';

import 'palette.dart';

@immutable
class FamilySpec {
  const FamilySpec({
    required this.light,
    required this.dark,
    required this.headingFont,
    required this.bodyFont,
    required this.numberFont,
    required this.radius,
    required this.radiusSmall,
  });

  final Palette light;
  final Palette dark;
  final String headingFont;
  final String bodyFont;
  final String numberFont;
  final double radius;
  final double radiusSmall;

  Palette palette(Brightness brightness) =>
      brightness == Brightness.dark ? dark : light;
}
