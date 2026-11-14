import 'dart:ui';

import 'package:material_ui/material_ui.dart';

import 'family_spec.dart';
import 'palette.dart';

@immutable
class LodestarTokens extends ThemeExtension<LodestarTokens> {
  const LodestarTokens({
    required this.palette,
    required this.headingFont,
    required this.bodyFont,
    required this.numberFont,
    required this.radius,
    required this.radiusSmall,
  });

  factory LodestarTokens.from(FamilySpec spec, Palette palette) {
    return LodestarTokens(
      palette: palette,
      headingFont: spec.headingFont,
      bodyFont: spec.bodyFont,
      numberFont: spec.numberFont,
      radius: spec.radius,
      radiusSmall: spec.radiusSmall,
    );
  }

  final Palette palette;
  final String headingFont;
  final String bodyFont;
  final String numberFont;
  final double radius;
  final double radiusSmall;

  BorderRadius get borderRadius => BorderRadius.circular(radius);

  BorderRadius get borderRadiusSmall => BorderRadius.circular(radiusSmall);

  @override
  LodestarTokens copyWith({
    Palette? palette,
    String? headingFont,
    String? bodyFont,
    String? numberFont,
    double? radius,
    double? radiusSmall,
  }) {
    return LodestarTokens(
      palette: palette ?? this.palette,
      headingFont: headingFont ?? this.headingFont,
      bodyFont: bodyFont ?? this.bodyFont,
      numberFont: numberFont ?? this.numberFont,
      radius: radius ?? this.radius,
      radiusSmall: radiusSmall ?? this.radiusSmall,
    );
  }

  @override
  LodestarTokens lerp(LodestarTokens? other, double t) {
    if (other == null) {
      return this;
    }
    final fontsFromOther = t >= 0.5;
    return LodestarTokens(
      palette: palette.lerp(other.palette, t),
      headingFont: fontsFromOther ? other.headingFont : headingFont,
      bodyFont: fontsFromOther ? other.bodyFont : bodyFont,
      numberFont: fontsFromOther ? other.numberFont : numberFont,
      radius: lerpDouble(radius, other.radius, t)!,
      radiusSmall: lerpDouble(radiusSmall, other.radiusSmall, t)!,
    );
  }
}
