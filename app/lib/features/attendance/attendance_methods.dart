class AttendanceMethods {
  const AttendanceMethods({required this.qr, required this.code});

  factory AttendanceMethods.fromJson(Map<String, Object?>? json) =>
      AttendanceMethods(
        qr: json?['qr'] as bool? ?? false,
        code: json?['code'] as bool? ?? false,
      );

  final bool qr;
  final bool code;
}
