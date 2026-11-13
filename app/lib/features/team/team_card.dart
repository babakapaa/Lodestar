import '../study/course_person.dart';
import 'team.dart';
import 'team_scope.dart';

class TeamCard {
  const TeamCard({
    required this.team,
    required this.scope,
    required this.invitedBy,
  });

  factory TeamCard.fromJson(Map<String, Object?> json) {
    final invitedBy = json['invitedBy'];
    return TeamCard(
      team: Team.fromJson(json['team']! as Map<String, Object?>),
      scope: TeamScope.fromJson(json['scope']! as Map<String, Object?>),
      invitedBy: invitedBy is Map<String, Object?>
          ? CoursePerson.fromJson(invitedBy)
          : null,
    );
  }

  final Team team;
  final TeamScope scope;
  final CoursePerson? invitedBy;
}
