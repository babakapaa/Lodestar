enum Grading {
  points,
  pass,
  none;

  static Grading parse(Object? value) => switch (value) {
    'points' => Grading.points,
    'pass' => Grading.pass,
    _ => Grading.none,
  };
}
