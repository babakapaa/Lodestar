import '../study/course_person.dart';

class DefenseSlot {
  const DefenseSlot({
    required this.id,
    required this.startsAt,
    required this.minutes,
    required this.place,
    required this.capacity,
    required this.taken,
    required this.mine,
    required this.students,
  });

  factory DefenseSlot.fromJson(Map<String, Object?> json) => DefenseSlot(
    id: json['id']! as String,
    startsAt: DateTime.parse(json['startsAt']! as String),
    minutes: (json['minutes']! as num).toInt(),
    place: json['place'] as String? ?? '',
    capacity: (json['capacity']! as num).toInt(),
    taken: (json['taken'] as num?)?.toInt() ?? 0,
    mine: json['mine'] as bool? ?? false,
    students: [
      for (final item in json['students'] as List<Object?>? ?? const [])
        CoursePerson.fromJson(item! as Map<String, Object?>),
    ],
  );

  final String id;
  final DateTime startsAt;
  final int minutes;
  final String place;
  final int capacity;
  final int taken;
  final bool mine;
  final List<CoursePerson> students;

  bool get free => taken < capacity;

  DateTime get endsAt => startsAt.add(Duration(minutes: minutes));
}
