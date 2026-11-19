class AuthorName {
  const AuthorName(this.family, this.initials);

  factory AuthorName.parse(String author) {
    final at = author.indexOf(',');
    if (at < 0) {
      return AuthorName(author.trim(), '');
    }
    final given = author.substring(at + 1).trim();
    final initials = [
      for (final part in given.split(RegExp(r'[\s.]+')))
        if (part.isNotEmpty)
          part.contains('-')
              ? part.split('-').map((p) => '${p[0]}.').join('-')
              : '${part[0]}.',
    ].join(' ');
    return AuthorName(author.substring(0, at).trim(), initials);
  }

  final String family;
  final String initials;
}
