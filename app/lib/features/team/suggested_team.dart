import '../study/course_person.dart';
import 'team_people.dart';

class SuggestedTeam {
  const SuggestedTeam({required this.members, required this.why});

  factory SuggestedTeam.fromJson(Map<String, Object?> json) => SuggestedTeam(
    members: TeamPeople.parse(json['members']),
    why: json['why'] as String? ?? '',
  );

  final List<CoursePerson> members;
  final String why;
}
