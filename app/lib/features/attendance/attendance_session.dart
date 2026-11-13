import 'attendance_counts.dart';
import 'attendance_json.dart';
import 'attendance_mark.dart';
import 'attendance_methods.dart';

class AttendanceSession {
  const AttendanceSession({
    required this.id,
    required this.seriesId,
    required this.date,
    required this.courseId,
    required this.courseTitle,
    required this.room,
    required this.startsAt,
    required this.endsAt,
    required this.open,
    required this.code,
    required this.codeUntil,
    required this.methods,
    required this.counts,
    required this.marks,
  });

  factory AttendanceSession.fromJson(Map<String, Object?> json) =>
      AttendanceSession(
        id: json['id']! as String,
        seriesId: json['seriesId']! as String,
        date: json['date']! as String,
        courseId: json['courseId']! as String,
        courseTitle: json['courseTitle'] as String? ?? '',
        room: json['room'] as String? ?? '',
        startsAt: AttendanceJson.time(json['startsAt']),
        endsAt: AttendanceJson.time(json['endsAt']),
        open: json['open'] as bool? ?? false,
        code: json['code'] as String? ?? '',
        codeUntil: AttendanceJson.maybeTime(json['codeUntil']),
        methods: AttendanceMethods.fromJson(
          json['methods'] as Map<String, Object?>?,
        ),
        counts: AttendanceCounts.fromJson(
          json['counts'] as Map<String, Object?>?,
        ),
        marks: AttendanceJson.list(json['marks'], AttendanceMark.fromJson),
      );

  final String id;
  final String seriesId;
  final String date;
  final String courseId;
  final String courseTitle;
  final String room;
  final DateTime startsAt;
  final DateTime endsAt;
  final bool open;
  final String code;
  final DateTime? codeUntil;
  final AttendanceMethods methods;
  final AttendanceCounts counts;
  final List<AttendanceMark> marks;
}
