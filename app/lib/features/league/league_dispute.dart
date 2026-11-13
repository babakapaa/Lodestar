import '../study/course_person.dart';
import 'dispute_status.dart';
import 'league_entry.dart';
import 'league_transition.dart';

class LeagueDispute {
  const LeagueDispute({
    required this.id,
    required this.entryId,
    required this.transitionId,
    required this.text,
    required this.status,
    required this.answer,
    required this.createdAt,
    this.person,
    this.entry,
    this.transition,
  });

  factory LeagueDispute.fromJson(Map<String, Object?> json) {
    final person = json['person'];
    final entry = json['entry'];
    final transition = json['transition'];
    return LeagueDispute(
      id: json['id']! as String,
      entryId: json['entryId'] as String?,
      transitionId: json['transitionId'] as String?,
      text: json['text'] as String? ?? '',
      status: DisputeStatus.parse(json['status']),
      answer: json['answer'] as String? ?? '',
      createdAt: DateTime.parse(json['createdAt']! as String),
      person: person is Map<String, Object?>
          ? CoursePerson.fromJson(person)
          : null,
      entry: entry is Map<String, Object?> ? LeagueEntry.fromJson(entry) : null,
      transition: transition is Map<String, Object?>
          ? LeagueTransition.fromJson(transition)
          : null,
    );
  }

  final String id;
  final String? entryId;
  final String? transitionId;
  final String text;
  final DisputeStatus status;
  final String answer;
  final DateTime createdAt;
  final CoursePerson? person;
  final LeagueEntry? entry;
  final LeagueTransition? transition;
}
