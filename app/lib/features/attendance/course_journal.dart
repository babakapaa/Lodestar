import 'attendance_json.dart';
import 'journal_lesson.dart';
import 'journal_student.dart';

class CourseJournal {
  const CourseJournal({
    required this.courseId,
    required this.courseTitle,
    required this.absentAlert,
    required this.lessons,
    required this.students,
  });

  factory CourseJournal.fromJson(Map<String, Object?> json) => CourseJournal(
    courseId: json['courseId']! as String,
    courseTitle: json['courseTitle'] as String? ?? '',
    absentAlert: json['absentAlert'] as int? ?? 3,
    lessons: AttendanceJson.list(json['lessons'], JournalLesson.fromJson),
    students: AttendanceJson.list(json['students'], JournalStudent.fromJson),
  );

  final String courseId;
  final String courseTitle;
  final int absentAlert;
  final List<JournalLesson> lessons;
  final List<JournalStudent> students;

  List<JournalStudent> get longAbsent => [
    for (final s in students)
      if (s.longAbsent) s,
  ];
}
