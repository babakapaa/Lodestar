abstract final class LeagueJson {
  static List<T> list<T>(
    Object? value,
    T Function(Map<String, Object?> json) parse,
  ) => [
    for (final item in value as List<Object?>? ?? const [])
      parse(item! as Map<String, Object?>),
  ];
}
