abstract final class NoteLinks {
  static final _pattern = RegExp(r'\[\[([^\[\]\n]+?)\]\]');

  static List<String> parse(String body) {
    final links = <String>{};
    for (final match in _pattern.allMatches(body)) {
      final target = match.group(1)!.split('|').first.split('#').first.trim();
      if (target.isNotEmpty) {
        links.add(target);
      }
    }
    return links.toList();
  }
}
