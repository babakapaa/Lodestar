import '../study/course_person.dart';

class TeamMember {
  const TeamMember({
    required this.person,
    required this.label,
    required this.joinedAt,
  });

  factory TeamMember.fromJson(Map<String, Object?> json) => TeamMember(
    person: CoursePerson.fromJson(json),
    label: json['label'] as String? ?? '',
    joinedAt: DateTime.parse(json['joinedAt']! as String),
  );

  final CoursePerson person;
  final String label;
  final DateTime joinedAt;
}
