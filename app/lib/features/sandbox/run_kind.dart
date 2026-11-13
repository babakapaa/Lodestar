enum RunKind {
  jupyter,
  vscode,
  job,
  demo;

  static RunKind parse(String? value) =>
      values.firstWhere((k) => k.name == value, orElse: () => job);

  bool get interactive => this == jupyter || this == vscode;

  bool get commanded => this == job || this == demo;
}
