import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';
import 'markdown_styles.dart';

class MarkdownController extends TextEditingController {
  MarkdownController({this.livePreview = false});

  final bool livePreview;
  bool _hide = false;

  TextStyle _marker(TextStyle base, MarkdownStyles styles) =>
      _hide ? styles.hidden : base.merge(styles.marker);

  static final _heading = RegExp(r'^(#{1,6})(\s)');
  static final _quote = RegExp(r'^(\s*>\s?)');
  static final _list = RegExp(r'^(\s*(?:[-*+]|\d+[.)])\s)(\[[ xX]\]\s)?');
  static final _inline = RegExp(
    r'(\[\[[^\[\]\n]+?\]\])|(\*\*[^*\n]+?\*\*)|(`[^`\n]+?`)|(==[^=\n]+?==)|(\*[^*\s][^*\n]*?\*)|(_[^_\s][^_\n]*?_)|((?<![\w#])#[\p{L}\p{N}_/-]+)',
    unicode: true,
  );

  @override
  TextSpan buildTextSpan({
    required BuildContext context,
    TextStyle? style,
    required bool withComposing,
  }) {
    if (withComposing && value.isComposingRangeValid) {
      return super.buildTextSpan(
        context: context,
        style: style,
        withComposing: withComposing,
      );
    }
    final styles = MarkdownStyles(style ?? const TextStyle(), context.palette);
    final lines = text.split('\n');
    final children = <InlineSpan>[];
    final selection = value.selection;
    final from = selection.isValid ? selection.start : -1;
    final to = selection.isValid ? selection.end : -1;
    var fenced = false;
    var offset = 0;
    for (var i = 0; i < lines.length; i++) {
      final line = lines[i];
      final end = offset + line.length;
      _hide = livePreview && !(from <= end && to >= offset);
      final fence = line.trimLeft().startsWith('```');
      if (fence || fenced) {
        children.add(TextSpan(text: line, style: styles.code));
        if (fence) {
          fenced = !fenced;
        }
      } else {
        children.addAll(_line(line, styles));
      }
      offset = end + 1;
      if (i < lines.length - 1) {
        children.add(TextSpan(text: '\n', style: style));
      }
    }
    return TextSpan(style: style, children: children);
  }

  List<InlineSpan> _line(String line, MarkdownStyles styles) {
    final heading = _heading.firstMatch(line);
    if (heading != null) {
      final level = heading.group(1)!.length;
      final base = styles.heading(level);
      final marker = heading.group(0)!;
      return [
        TextSpan(text: marker, style: _marker(base, styles)),
        ..._inlines(line.substring(marker.length), base, styles),
      ];
    }
    final quote = _quote.firstMatch(line);
    if (quote != null) {
      final marker = quote.group(0)!;
      return [
        TextSpan(
          text: marker,
          style: _hide ? styles.hidden : styles.accentMarker,
        ),
        ..._inlines(line.substring(marker.length), styles.quote, styles),
      ];
    }
    final list = _list.firstMatch(line);
    if (list != null) {
      final bullet = list.group(1)!;
      final task = list.group(2);
      final done = task != null && task.toLowerCase().contains('x');
      final rest = line.substring(list.group(0)!.length);
      return [
        TextSpan(text: bullet, style: styles.accentMarker),
        if (task != null) TextSpan(text: task, style: styles.accentMarker),
        ..._inlines(rest, done ? styles.done : styles.base, styles),
      ];
    }
    return _inlines(line, styles.base, styles);
  }

  List<InlineSpan> _inlines(
    String text,
    TextStyle base,
    MarkdownStyles styles,
  ) {
    final spans = <InlineSpan>[];
    var at = 0;
    for (final match in _inline.allMatches(text)) {
      if (match.start > at) {
        spans.add(TextSpan(text: text.substring(at, match.start), style: base));
      }
      final token = match.group(0)!;
      spans.addAll(switch (match) {
        _ when match.group(1) != null => _wrapped(
          token,
          2,
          base,
          styles.link,
          styles,
        ),
        _ when match.group(2) != null => _wrapped(
          token,
          2,
          base,
          styles.bold,
          styles,
        ),
        _ when match.group(3) != null => _wrapped(
          token,
          1,
          base,
          styles.inlineCode,
          styles,
        ),
        _ when match.group(4) != null => _wrapped(
          token,
          2,
          base,
          styles.highlight,
          styles,
        ),
        _ when match.group(5) != null || match.group(6) != null => _wrapped(
          token,
          1,
          base,
          styles.italic,
          styles,
        ),
        _ => [TextSpan(text: token, style: base.merge(styles.tag))],
      });
      at = match.end;
    }
    if (at < text.length) {
      spans.add(TextSpan(text: text.substring(at), style: base));
    }
    return spans;
  }

  List<InlineSpan> _wrapped(
    String token,
    int width,
    TextStyle base,
    TextStyle inner,
    MarkdownStyles styles,
  ) {
    final marker = _marker(base, styles);
    return [
      TextSpan(text: token.substring(0, width), style: marker),
      TextSpan(
        text: token.substring(width, token.length - width),
        style: base.merge(inner),
      ),
      TextSpan(text: token.substring(token.length - width), style: marker),
    ];
  }
}
