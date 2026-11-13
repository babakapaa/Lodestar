import '../study/course_person.dart';
import 'history_event.dart';

class HistoryEntry {
  const HistoryEntry({
    required this.person,
    required this.event,
    required this.actor,
    required this.at,
  });

  factory HistoryEntry.fromJson(Map<String, Object?> json) {
    final person = json['person'];
    final actor = json['actor'];
    return HistoryEntry(
      person: person is Map<String, Object?>
          ? CoursePerson.fromJson(person)
          : null,
      event: HistoryEvent.parse(json['event']),
      actor: actor is Map<String, Object?>
          ? CoursePerson.fromJson(actor)
          : null,
      at: DateTime.parse(json['at']! as String),
    );
  }

  final CoursePerson? person;
  final HistoryEvent event;
  final CoursePerson? actor;
  final DateTime at;
}
