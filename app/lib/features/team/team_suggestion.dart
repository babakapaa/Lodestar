import '../study/course_person.dart';
import 'suggested_team.dart';
import 'team_people.dart';

class TeamSuggestion {
  const TeamSuggestion({required this.teams, required this.without});

  factory TeamSuggestion.fromJson(Map<String, Object?> json) => TeamSuggestion(
    teams: [
      for (final item in json['teams'] as List<Object?>? ?? const [])
        SuggestedTeam.fromJson(item! as Map<String, Object?>),
    ],
    without: TeamPeople.parse(json['without']),
  );

  final List<SuggestedTeam> teams;
  final List<CoursePerson> without;
}
