import 'md_block.dart';

abstract final class MdParser {
  static final _heading = RegExp(r'^(#{1,6})\s+(.*)$');
  static final _item = RegExp(r'^(\s*)([-*+]|\d+[.)])\s+(\[[ xX]\]\s+)?(.*)$');
  static final _quote = RegExp(r'^\s*>\s?(.*)$');
  static final _rule = RegExp(r'^\s*([-*_])(\s*\1){2,}\s*$');

  static List<MdBlock> parse(String source) {
    final lines = source.replaceAll('\r\n', '\n').split('\n');
    final blocks = <MdBlock>[];
    final paragraph = <String>[];

    void flush() {
      if (paragraph.isNotEmpty) {
        blocks.add(MdBlock(MdKind.paragraph, paragraph.join('\n')));
        paragraph.clear();
      }
    }

    var i = 0;
    while (i < lines.length) {
      final line = lines[i];
      final trimmed = line.trim();
      if (trimmed.startsWith('```')) {
        flush();
        final code = <String>[];
        i++;
        while (i < lines.length && !lines[i].trim().startsWith('```')) {
          code.add(lines[i]);
          i++;
        }
        blocks.add(MdBlock(MdKind.code, code.join('\n')));
        i++;
        continue;
      }
      final open = trimmed.startsWith(r'$$')
          ? r'$$'
          : trimmed.startsWith(r'\[')
          ? r'\['
          : null;
      if (open != null) {
        flush();
        final close = open == r'$$' ? r'$$' : r'\]';
        var body = trimmed.substring(2);
        if (body.trimRight().endsWith(close) && body.trim().length >= 2) {
          body = body.trimRight();
          blocks.add(
            MdBlock(MdKind.math, body.substring(0, body.length - 2).trim()),
          );
          i++;
          continue;
        }
        final math = <String>[if (body.trim().isNotEmpty) body];
        i++;
        while (i < lines.length && !lines[i].trimRight().endsWith(close)) {
          math.add(lines[i]);
          i++;
        }
        if (i < lines.length) {
          final last = lines[i].trimRight();
          math.add(last.substring(0, last.length - 2));
        }
        blocks.add(MdBlock(MdKind.math, math.join('\n').trim()));
        i++;
        continue;
      }
      if (trimmed.isEmpty) {
        flush();
        i++;
        continue;
      }
      if (_rule.hasMatch(line)) {
        flush();
        blocks.add(const MdBlock(MdKind.rule, ''));
        i++;
        continue;
      }
      final heading = _heading.firstMatch(line);
      if (heading != null) {
        flush();
        blocks.add(
          MdBlock(
            MdKind.heading,
            heading.group(2)!.trim(),
            level: heading.group(1)!.length,
          ),
        );
        i++;
        continue;
      }
      final item = _item.firstMatch(line);
      if (item != null) {
        flush();
        final task = item.group(3);
        final marker = item.group(2)!;
        blocks.add(
          MdBlock(
            MdKind.item,
            item.group(4)!,
            level: item.group(1)!.length ~/ 2,
            marker: RegExp(r'\d').hasMatch(marker) ? marker : '•',
            done: task?.toLowerCase().contains('x'),
          ),
        );
        i++;
        continue;
      }
      final quote = _quote.firstMatch(line);
      if (quote != null) {
        flush();
        final text = <String>[quote.group(1)!];
        i++;
        while (i < lines.length) {
          final next = _quote.firstMatch(lines[i]);
          if (next == null) {
            break;
          }
          text.add(next.group(1)!);
          i++;
        }
        blocks.add(MdBlock(MdKind.quote, text.join('\n')));
        continue;
      }
      paragraph.add(line.trimRight());
      i++;
    }
    flush();
    return blocks;
  }
}
