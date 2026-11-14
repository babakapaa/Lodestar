import 'package:material_ui/material_ui.dart';

@immutable
class Palette {
  const Palette({
    required this.bg,
    required this.side,
    required this.surface,
    required this.surface2,
    required this.border,
    required this.text,
    required this.muted,
    required this.accent,
    required this.onAccent,
    required this.accentSoft,
    required this.warn,
    required this.warnSoft,
    required this.shadow,
  });

  final Color bg;
  final Color side;
  final Color surface;
  final Color surface2;
  final Color border;
  final Color text;
  final Color muted;
  final Color accent;
  final Color onAccent;
  final Color accentSoft;
  final Color warn;
  final Color warnSoft;
  final List<BoxShadow> shadow;

  Palette withAccent({
    required Color accent,
    required Color onAccent,
    required Color accentSoft,
  }) {
    return Palette(
      bg: bg,
      side: side,
      surface: surface,
      surface2: surface2,
      border: border,
      text: text,
      muted: muted,
      accent: accent,
      onAccent: onAccent,
      accentSoft: accentSoft,
      warn: warn,
      warnSoft: warnSoft,
      shadow: shadow,
    );
  }

  Palette lerp(Palette other, double t) {
    Color mix(Color a, Color b) => Color.lerp(a, b, t)!;
    return Palette(
      bg: mix(bg, other.bg),
      side: mix(side, other.side),
      surface: mix(surface, other.surface),
      surface2: mix(surface2, other.surface2),
      border: mix(border, other.border),
      text: mix(text, other.text),
      muted: mix(muted, other.muted),
      accent: mix(accent, other.accent),
      onAccent: mix(onAccent, other.onAccent),
      accentSoft: mix(accentSoft, other.accentSoft),
      warn: mix(warn, other.warn),
      warnSoft: mix(warnSoft, other.warnSoft),
      shadow: BoxShadow.lerpList(shadow, other.shadow, t) ?? other.shadow,
    );
  }
}
