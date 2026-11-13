enum NavigatorMode {
  course,
  study,
  project;

  static NavigatorMode parse(Object? value) => switch (value) {
    'study' => study,
    'project' => project,
    _ => course,
  };

  String? get wire => this == course ? null : name;
}
