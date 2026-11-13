import '../study/course_person.dart';
import 'leave_request.dart';
import 'team.dart';
import 'team_people.dart';
import 'team_scope.dart';
import 'team_seeker.dart';

class TeamScopeView {
  const TeamScopeView({
    required this.scope,
    required this.leads,
    required this.myTeam,
    required this.solo,
    required this.seeking,
    required this.seekers,
    required this.teams,
    required this.unplaced,
    required this.solos,
    required this.requests,
  });

  factory TeamScopeView.fromJson(Map<String, Object?> json) {
    final myTeam = json['myTeam'];
    return TeamScopeView(
      scope: TeamScope.fromJson(json['scope']! as Map<String, Object?>),
      leads: json['role'] == 'lead',
      myTeam: myTeam is Map<String, Object?> ? Team.fromJson(myTeam) : null,
      solo: json['solo'] as bool? ?? false,
      seeking: json['seeking'] as bool? ?? false,
      seekers: [
        for (final item in json['seekers'] as List<Object?>? ?? const [])
          TeamSeeker.fromJson(item! as Map<String, Object?>),
      ],
      teams: [
        for (final item in json['teams'] as List<Object?>? ?? const [])
          Team.fromJson(item! as Map<String, Object?>),
      ],
      unplaced: TeamPeople.parse(json['unplaced']),
      solos: TeamPeople.parse(json['solos']),
      requests: [
        for (final item in json['requests'] as List<Object?>? ?? const [])
          LeaveRequest.fromJson(item! as Map<String, Object?>),
      ],
    );
  }

  final TeamScope scope;
  final bool leads;
  final Team? myTeam;
  final bool solo;
  final bool seeking;
  final List<TeamSeeker> seekers;
  final List<Team> teams;
  final List<CoursePerson> unplaced;
  final List<CoursePerson> solos;
  final List<LeaveRequest> requests;
}
