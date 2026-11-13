import '../study/course_person.dart';

class TeamSeeker {
  const TeamSeeker({required this.person, required this.note});

  factory TeamSeeker.fromJson(Map<String, Object?> json) => TeamSeeker(
    person: CoursePerson.fromJson(json),
    note: json['note'] as String? ?? '',
  );

  final CoursePerson person;
  final String note;
}
