import 'attendance_status.dart';

class AttendanceMark {
  const AttendanceMark({
    required this.userId,
    required this.name,
    required this.status,
    required this.method,
  });

  factory AttendanceMark.fromJson(Map<String, Object?> json) => AttendanceMark(
    userId: json['userId']! as String,
    name: json['name'] as String? ?? '',
    status: AttendanceStatus.parse(json['status']),
    method: json['method'] as String?,
  );

  final String userId;
  final String name;
  final AttendanceStatus? status;
  final String? method;
}
