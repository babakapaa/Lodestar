import 'attendance_json.dart';

class AttendanceQr {
  const AttendanceQr({required this.token, required this.expiresAt});

  factory AttendanceQr.fromJson(Map<String, Object?> json) => AttendanceQr(
    token: json['token']! as String,
    expiresAt: AttendanceJson.time(json['expiresAt']),
  );

  final String token;
  final DateTime expiresAt;
}
