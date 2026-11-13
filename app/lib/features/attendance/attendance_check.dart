import 'attendance_json.dart';
import 'attendance_status.dart';

class AttendanceCheck {
  const AttendanceCheck({
    required this.courseTitle,
    required this.startsAt,
    required this.status,
  });

  factory AttendanceCheck.fromJson(Map<String, Object?> json) =>
      AttendanceCheck(
        courseTitle: json['courseTitle'] as String? ?? '',
        startsAt: AttendanceJson.time(json['startsAt']),
        status:
            AttendanceStatus.parse(json['status']) ?? AttendanceStatus.present,
      );

  final String courseTitle;
  final DateTime startsAt;
  final AttendanceStatus status;
}
