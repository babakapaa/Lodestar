import 'package:material_ui/material_ui.dart';

import 'palette.dart';

enum AccentChoice {
  auto(null, null),
  blue(Color(0xFF1D4ED8), Color(0xFF93C5FD)),
  indigo(Color(0xFF4338CA), Color(0xFFA5B4FC)),
  violet(Color(0xFF6D28D9), Color(0xFFC4B5FD)),
  pink(Color(0xFFBE185D), Color(0xFFF9A8D4)),
  red(Color(0xFFB91C1C), Color(0xFFFCA5A5)),
  amber(Color(0xFF92400E), Color(0xFFFBBF24)),
  green(Color(0xFF166534), Color(0xFF4ADE80)),
  teal(Color(0xFF0F766E), Color(0xFF2DD4BF));

  const AccentChoice(this.light, this.dark);

  static const _onLight = Color(0xFFFFFFFF);
  static const _onDark = Color(0xFF111111);

  final Color? light;
  final Color? dark;

  Palette apply(Palette palette, Brightness brightness) {
    final accent = brightness == Brightness.dark ? dark : light;
    if (accent == null) {
      return palette;
    }
    final soft = Color.alphaBlend(
      accent.withValues(alpha: brightness == Brightness.dark ? 0.14 : 0.1),
      palette.surface,
    );
    return palette.withAccent(
      accent: accent,
      onAccent: brightness == Brightness.dark ? _onDark : _onLight,
      accentSoft: soft,
    );
  }
}
