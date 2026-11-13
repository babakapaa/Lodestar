abstract final class SandboxJson {
  static List<T> list<T>(
    Object? value,
    T Function(Map<String, Object?> json) parse,
  ) => [
    for (final item in value as List<Object?>? ?? const [])
      parse(item! as Map<String, Object?>),
  ];

  static double number(Object? value) => (value as num?)?.toDouble() ?? 0;

  static double? maybeNumber(Object? value) => (value as num?)?.toDouble();

  static DateTime? time(Object? value) =>
      value == null ? null : DateTime.parse(value as String).toLocal();
}
