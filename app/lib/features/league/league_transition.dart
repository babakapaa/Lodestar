import '../study/course_person.dart';
import 'transition_status.dart';

class LeagueTransition {
  const LeagueTransition({
    required this.id,
    required this.from,
    required this.to,
    required this.points,
    required this.status,
    required this.createdAt,
    this.person,
  });

  factory LeagueTransition.fromJson(Map<String, Object?> json) {
    final person = json['person'];
    return LeagueTransition(
      id: json['id']! as String,
      from: (json['from'] as num?)?.toInt() ?? 0,
      to: (json['to'] as num?)?.toInt() ?? 0,
      points: (json['points'] as num?)?.toInt() ?? 0,
      status: TransitionStatus.parse(json['status']),
      createdAt: DateTime.parse(json['createdAt']! as String),
      person: person is Map<String, Object?>
          ? CoursePerson.fromJson(person)
          : null,
    );
  }

  final String id;
  final int from;
  final int to;
  final int points;
  final TransitionStatus status;
  final DateTime createdAt;
  final CoursePerson? person;

  bool get up => to > from;
}
