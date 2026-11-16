import 'package:flutter_math_fork/flutter_math.dart';
import 'package:material_ui/material_ui.dart';

import '../tokens/palette.dart';
import '../typography/font_families.dart';

class MdInline {
  MdInline(this.base, this.palette, {this.onLink});

  static final _pattern = RegExp(
    r'(\\\((.+?)\\\))|(\$([^$\n]+?)\$)|(\*\*([^*\n]+?)\*\*)|(__([^_\n]+?)__)|(`([^`\n]+?)`)|(\[\[([^\[\]\n]+?)\]\])|(\[([^\]\n]+?)\]\(([^)\s]+)\))|(==([^=\n]+?)==)|(\*([^*\s][^*\n]*?)\*)|((?<![\p{L}\p{N}])_([^_\s][^_\n]*?)_(?![\p{L}\p{N}]))',
    unicode: true,
  );

  final TextStyle base;
  final Palette palette;
  final ValueChanged<String>? onLink;

  TextStyle get _link => base.copyWith(
    color: palette.accent,
    decoration: TextDecoration.underline,
    decorationColor: palette.accent.withValues(alpha: 0.4),
  );

  List<InlineSpan> spans(String text, [TextStyle? style]) {
    final current = style ?? base;
    final out = <InlineSpan>[];
    var at = 0;
    for (final m in _pattern.allMatches(text)) {
      if (m.start > at) {
        out.add(TextSpan(text: text.substring(at, m.start), style: current));
      }
      if (m.group(1) != null || m.group(3) != null) {
        out.add(math(m.group(2) ?? m.group(4)!, current));
      } else if (m.group(5) != null || m.group(7) != null) {
        out.addAll(
          spans(
            m.group(6) ?? m.group(8)!,
            current.copyWith(fontWeight: FontWeight.w700),
          ),
        );
      } else if (m.group(9) != null) {
        out.add(
          TextSpan(
            text: m.group(10),
            style: current.copyWith(
              fontFamily: FontFamilies.jetBrainsMono,
              fontSize: (current.fontSize ?? 15) * 0.9,
              backgroundColor: palette.surface2,
            ),
          ),
        );
      } else if (m.group(11) != null) {
        final target = m.group(12)!;
        final shown = target.contains('|') ? target.split('|').last : target;
        out.add(TextSpan(text: shown, style: _link));
      } else if (m.group(13) != null) {
        out.add(TextSpan(text: m.group(14), style: _link));
      } else if (m.group(16) != null) {
        out.add(
          TextSpan(
            text: m.group(17),
            style: current.copyWith(backgroundColor: palette.warnSoft),
          ),
        );
      } else {
        out.addAll(
          spans(
            m.group(19) ?? m.group(21)!,
            current.copyWith(fontStyle: FontStyle.italic),
          ),
        );
      }
      at = m.end;
    }
    if (at < text.length) {
      out.add(TextSpan(text: text.substring(at), style: current));
    }
    return out;
  }

  InlineSpan math(String tex, TextStyle style) {
    return WidgetSpan(
      alignment: PlaceholderAlignment.middle,
      child: Math.tex(
        tex,
        mathStyle: MathStyle.text,
        textStyle: style,
        onErrorFallback: (_) => Text(tex, style: style),
      ),
    );
  }
}
