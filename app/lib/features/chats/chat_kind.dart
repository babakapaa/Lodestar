enum ChatKind {
  direct,
  group,
  channel;

  static ChatKind parse(String value) => values.byName(value);
}
