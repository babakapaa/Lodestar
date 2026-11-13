import 'attendance_counts.dart';
import 'attendance_json.dart';

class JournalLesson {
  const JournalLesson({
    required this.seriesId,
    required this.date,
    required this.startsAt,
    required this.room,
    required this.cancelled,
    required this.sessionId,
    required this.open,
    required this.counts,
  });

  factory JournalLesson.fromJson(Map<String, Object?> json) => JournalLesson(
    seriesId: json['seriesId']! as String,
    date: json['date']! as String,
    startsAt: AttendanceJson.time(json['startsAt']),
    room: json['room'] as String? ?? '',
    cancelled: json['cancelled'] as bool? ?? false,
    sessionId: json['sessionId'] as String?,
    open: json['open'] as bool? ?? false,
    counts: AttendanceCounts.fromJson(json['counts'] as Map<String, Object?>?),
  );

  final String seriesId;
  final String date;
  final DateTime startsAt;
  final String room;
  final bool cancelled;
  final String? sessionId;
  final bool open;
  final AttendanceCounts counts;
}
