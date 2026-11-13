import 'attendance_json.dart';
import 'attendance_methods.dart';
import 'attendance_status.dart';

class OpenLesson {
  const OpenLesson({
    required this.sessionId,
    required this.courseTitle,
    required this.room,
    required this.startsAt,
    required this.endsAt,
    required this.status,
    required this.methods,
  });

  factory OpenLesson.fromJson(Map<String, Object?> json) => OpenLesson(
    sessionId: json['sessionId']! as String,
    courseTitle: json['courseTitle'] as String? ?? '',
    room: json['room'] as String? ?? '',
    startsAt: AttendanceJson.time(json['startsAt']),
    endsAt: AttendanceJson.time(json['endsAt']),
    status: AttendanceStatus.parse(json['status']),
    methods: AttendanceMethods.fromJson(
      json['methods'] as Map<String, Object?>?,
    ),
  );

  final String sessionId;
  final String courseTitle;
  final String room;
  final DateTime startsAt;
  final DateTime endsAt;
  final AttendanceStatus? status;
  final AttendanceMethods methods;
}
