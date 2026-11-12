import 'package:material_ui/material_ui.dart';

import '../tokens/family_spec.dart';
import '../tokens/palette.dart';
import 'font_families.dart';

class LodestarTypography {
  const LodestarTypography(this.spec, this.palette);

  final FamilySpec spec;
  final Palette palette;

  TextTheme build() {
    return TextTheme(
      displaySmall: _heading(30, FontWeight.w600, 1.2),
      headlineMedium: _heading(26, FontWeight.w600, 1.25),
      headlineSmall: _heading(22, FontWeight.w600, 1.3),
      titleLarge: _heading(18, FontWeight.w600, 1.35),
      titleMedium: _body(16, FontWeight.w600, 1.4, palette.text),
      titleSmall: _body(14, FontWeight.w600, 1.4, palette.text),
      bodyLarge: _body(16, FontWeight.w400, 1.5, palette.text),
      bodyMedium: _body(14, FontWeight.w400, 1.5, palette.text),
      bodySmall: _body(12, FontWeight.w400, 1.45, palette.muted),
      labelLarge: _body(14, FontWeight.w600, 1.3, palette.text),
      labelMedium: _body(13, FontWeight.w500, 1.3, palette.text),
      labelSmall: _body(12, FontWeight.w500, 1.3, palette.muted),
    );
  }

  TextStyle numbers(double size, FontWeight weight) =>
      _style(spec.numberFont, size, weight, 1.3, palette.text);

  TextStyle _heading(double size, FontWeight weight, double height) =>
      _style(spec.headingFont, size, weight, height, palette.text);

  TextStyle _body(double size, FontWeight weight, double height, Color color) =>
      _style(spec.bodyFont, size, weight, height, color);

  TextStyle _style(
    String family,
    double size,
    FontWeight weight,
    double height,
    Color color,
  ) {
    return TextStyle(
      fontFamily: family,
      fontFamilyFallback: FontFamilies.fallback,
      fontSize: size,
      fontWeight: weight,
      fontVariations: [FontVariation.weight(weight.value.toDouble())],
      height: height,
      color: color,
    );
  }
}
