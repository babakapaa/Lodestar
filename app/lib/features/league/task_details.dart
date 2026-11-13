import '../study/course_person.dart';
import '../team/team_people.dart';
import 'league_solution.dart';
import 'league_task.dart';

class TaskDetails {
  const TaskDetails({
    required this.task,
    required this.league,
    required this.kind,
    required this.originId,
    required this.leads,
    required this.solutions,
    required this.people,
  });

  factory TaskDetails.fromJson(Map<String, Object?> json) => TaskDetails(
    task: LeagueTask.fromJson(json['task']! as Map<String, Object?>),
    league: json['league'] as String? ?? '',
    kind: json['kind'] as String? ?? '',
    originId: json['originId'] as String? ?? '',
    leads: json['role'] == 'lead',
    solutions: [
      for (final item in json['solutions'] as List<Object?>? ?? const [])
        LeagueSolution.fromJson(item! as Map<String, Object?>),
    ],
    people: TeamPeople.parse(json['people']),
  );

  final LeagueTask task;
  final String league;
  final String kind;
  final String originId;
  final bool leads;
  final List<LeagueSolution> solutions;
  final List<CoursePerson> people;

  String nameOf(String userId) {
    for (final person in people) {
      if (person.userId == userId) {
        return person.label;
      }
    }
    return userId;
  }
}
