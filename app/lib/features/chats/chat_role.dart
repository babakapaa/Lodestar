enum ChatRole {
  owner,
  admin,
  member;

  static ChatRole parse(String value) => values.byName(value);

  bool get manages => this == owner || this == admin;
}
