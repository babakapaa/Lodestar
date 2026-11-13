enum DefenseFormat {
  none,
  full,
  short,
  sample;

  static DefenseFormat parse(Object? value) => switch (value) {
    'full' => full,
    'short' => short,
    'sample' => sample,
    _ => none,
  };

  String? get wire => this == none ? null : name;
}
