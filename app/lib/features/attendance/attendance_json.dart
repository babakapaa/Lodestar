abstract final class AttendanceJson {
  static List<T> list<T>(
    Object? value,
    T Function(Map<String, Object?> json) parse,
  ) => [
    for (final item in value as List<Object?>? ?? const [])
      parse(item! as Map<String, Object?>),
  ];

  static DateTime time(Object? value) =>
      DateTime.parse(value! as String).toLocal();

  static DateTime? maybeTime(Object? value) =>
      value == null ? null : DateTime.parse(value as String).toLocal();
}
