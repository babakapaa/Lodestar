enum NoteRole {
  owner,
  editor,
  reader;

  bool get canEdit => this != reader;

  static NoteRole parse(String name) =>
      values.firstWhere((role) => role.name == name, orElse: () => reader);
}
