class FrontMatter {
  const FrontMatter(this.tags, this.body);

  final List<String> tags;
  final String body;

  static final _block = RegExp(r'^---\r?\n([\s\S]*?)\r?\n---[ \t]*(?:\r?\n|$)');
  static final _inline = RegExp(r'^tags:\s*(.*)$');
  static final _item = RegExp(r'^\s+-\s*(.+)$');

  static FrontMatter parse(String text) {
    final match = _block.firstMatch(text);
    if (match == null) {
      return FrontMatter(const [], text);
    }
    final kept = <String>[];
    final tags = <String>[];
    var inTags = false;
    for (final line in match.group(1)!.split(RegExp(r'\r?\n'))) {
      final inline = _inline.firstMatch(line);
      if (inline != null) {
        inTags = true;
        tags.addAll(_split(inline.group(1)!));
        continue;
      }
      final item = inTags ? _item.firstMatch(line) : null;
      if (item != null) {
        tags.addAll(_split(item.group(1)!));
        continue;
      }
      inTags = false;
      kept.add(line);
    }
    final rest = text.substring(match.end);
    final body = kept.any((line) => line.trim().isNotEmpty)
        ? '---\n${kept.join('\n')}\n---\n$rest'
        : rest;
    return FrontMatter(tags.toSet().toList(), body);
  }

  static String render(List<String> tags, String body) {
    if (tags.isEmpty) {
      return body;
    }
    final lines = ['tags:', for (final tag in tags) '  - $tag'].join('\n');
    final match = _block.firstMatch(body);
    if (match == null) {
      return '---\n$lines\n---\n$body';
    }
    return '---\n$lines\n${match.group(1)}\n---\n${body.substring(match.end)}';
  }

  static List<String> _split(String value) => [
    for (final part
        in value
            .replaceAll(
              RegExp(
                r'[\[\]"'
                "'"
                r']',
              ),
              '',
            )
            .split(RegExp(r'[,\s]+')))
      if (part.replaceAll('#', '').isNotEmpty) part.replaceAll('#', ''),
  ];
}
