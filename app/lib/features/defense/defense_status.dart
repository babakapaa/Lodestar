enum DefenseStatus {
  planned,
  held,
  missed;

  static DefenseStatus parse(Object? value) => switch (value) {
    'held' => held,
    'missed' => missed,
    _ => planned,
  };

  bool get finished => this != planned;
}
