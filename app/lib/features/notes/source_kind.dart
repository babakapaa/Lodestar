enum SourceKind {
  book,
  article,
  chapter,
  web,
  thesis,
  other;

  static SourceKind parse(String name) =>
      values.firstWhere((kind) => kind.name == name, orElse: () => other);
}
