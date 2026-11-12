import 'deadline.dart';
import 'lesson.dart';

class Agenda {
  const Agenda({
    required this.date,
    required this.lessons,
    required this.deadlines,
  });

  factory Agenda.fromJson(Map<String, Object?> json) {
    return Agenda(
      date: json['date']! as String,
      lessons: [
        for (final item in json['lessons'] as List<Object?>? ?? const [])
          Lesson.fromJson(item! as Map<String, Object?>),
      ],
      deadlines: [
        for (final item in json['deadlines'] as List<Object?>? ?? const [])
          Deadline.fromJson(item! as Map<String, Object?>),
      ],
    );
  }

  final String date;
  final List<Lesson> lessons;
  final List<Deadline> deadlines;
}
