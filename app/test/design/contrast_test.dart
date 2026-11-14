import 'dart:math';

import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/design/tokens/accent_choice.dart';
import 'package:lodestar/design/tokens/theme_family.dart';
import 'package:material_ui/material_ui.dart';

double contrast(Color a, Color b) {
  final la = a.computeLuminance();
  final lb = b.computeLuminance();
  return (max(la, lb) + 0.05) / (min(la, lb) + 0.05);
}

void main() {
  for (final family in ThemeFamily.values) {
    for (final brightness in Brightness.values) {
      for (final accent in AccentChoice.values) {
        final p = accent.apply(family.spec.palette(brightness), brightness);
        final name = '${family.name} ${brightness.name} ${accent.name}';

        test('$name: text is readable on every background', () {
          for (final bg in [p.bg, p.side, p.surface, p.surface2]) {
            expect(contrast(p.text, bg), greaterThanOrEqualTo(4.5));
          }
        });

        test('$name: muted text is readable on main backgrounds', () {
          for (final bg in [p.bg, p.surface]) {
            expect(contrast(p.muted, bg), greaterThanOrEqualTo(4.5));
          }
        });

        test('$name: accent button label is readable', () {
          expect(contrast(p.onAccent, p.accent), greaterThanOrEqualTo(4.5));
        });

        test('$name: accent text is readable on soft accent', () {
          expect(contrast(p.accent, p.accentSoft), greaterThanOrEqualTo(4.5));
        });
      }
    }
  }
}
