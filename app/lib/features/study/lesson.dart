import 'lesson_status.dart';

class Lesson {
  const Lesson({
    required this.seriesId,
    required this.courseId,
    required this.courseTitle,
    required this.kind,
    required this.date,
    required this.startsAt,
    required this.endsAt,
    required this.room,
    required this.status,
    required this.note,
  });

  factory Lesson.fromJson(Map<String, Object?> json) {
    return Lesson(
      seriesId: json['seriesId']! as String,
      courseId: json['courseId']! as String,
      courseTitle: json['courseTitle'] as String? ?? '',
      kind: json['kind'] as String? ?? '',
      date: json['date']! as String,
      startsAt: DateTime.parse(json['startsAt']! as String),
      endsAt: DateTime.parse(json['endsAt']! as String),
      room: json['room'] as String? ?? '',
      status: LessonStatus.parse(json['status']),
      note: json['note'] as String? ?? '',
    );
  }

  final String seriesId;
  final String courseId;
  final String courseTitle;
  final String kind;
  final String date;
  final DateTime startsAt;
  final DateTime endsAt;
  final String room;
  final LessonStatus status;
  final String note;
}
