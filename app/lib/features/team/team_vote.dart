import '../study/course_person.dart';
import 'vote_decision.dart';

class TeamVote {
  const TeamVote({
    required this.person,
    required this.decision,
    required this.comment,
  });

  factory TeamVote.fromJson(Map<String, Object?> json) => TeamVote(
    person: CoursePerson.fromJson(json),
    decision: VoteDecision.parse(json['decision']),
    comment: json['comment'] as String? ?? '',
  );

  final CoursePerson person;
  final VoteDecision decision;
  final String comment;
}
