import 'package:material_ui/material_ui.dart';

import '../../design/tokens/palette.dart';
import '../../design/typography/font_families.dart';

class MarkdownStyles {
  MarkdownStyles(this.base, this.palette);

  static const _scales = [1.6, 1.4, 1.25, 1.12, 1.05, 1.0];

  final TextStyle base;
  final Palette palette;

  double get _size => base.fontSize ?? 16;

  TextStyle heading(int level) => base.copyWith(
    fontSize: _size * _scales[(level - 1).clamp(0, _scales.length - 1)],
    fontWeight: FontWeight.w700,
    height: 1.4,
  );

  TextStyle get marker =>
      TextStyle(color: palette.muted.withValues(alpha: 0.6));

  TextStyle get accentMarker => base.copyWith(color: palette.accent);

  TextStyle get quote =>
      base.copyWith(color: palette.muted, fontStyle: FontStyle.italic);

  TextStyle get done => base.copyWith(
    color: palette.muted,
    decoration: TextDecoration.lineThrough,
  );

  TextStyle get bold => const TextStyle(fontWeight: FontWeight.w700);

  TextStyle get italic => const TextStyle(fontStyle: FontStyle.italic);

  TextStyle get link => TextStyle(
    color: palette.accent,
    decoration: TextDecoration.underline,
    decorationColor: palette.accent.withValues(alpha: 0.4),
  );

  TextStyle get tag =>
      TextStyle(color: palette.accent, backgroundColor: palette.accentSoft);

  TextStyle get highlight => TextStyle(backgroundColor: palette.warnSoft);

  TextStyle get inlineCode => TextStyle(
    fontFamily: FontFamilies.jetBrainsMono,
    fontSize: _size * 0.9,
    backgroundColor: palette.surface2,
  );

  TextStyle get code => base.merge(inlineCode);
}
