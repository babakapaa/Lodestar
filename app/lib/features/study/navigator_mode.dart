enum NavigatorMode {
  course,
  study,
  project,
  closed;

  static NavigatorMode parse(Object? value) => switch (value) {
    'study' => study,
    'project' => project,
    'closed' => closed,
    _ => course,
  };

  String? get wire => this == course ? null : name;
}
