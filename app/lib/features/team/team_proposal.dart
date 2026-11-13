import '../study/course_person.dart';
import 'proposal_status.dart';
import 'team_vote.dart';
import 'team_work.dart';
import 'vote_decision.dart';

class TeamProposal {
  const TeamProposal({
    required this.id,
    required this.proposedBy,
    required this.work,
    required this.note,
    required this.status,
    required this.silentAt,
    required this.votes,
  });

  factory TeamProposal.fromJson(Map<String, Object?> json) {
    final proposedBy = json['proposedBy'];
    return TeamProposal(
      id: json['id']! as String,
      proposedBy: proposedBy is Map<String, Object?>
          ? CoursePerson.fromJson(proposedBy)
          : null,
      work: TeamWork.fromJson(json['work']),
      note: json['note'] as String? ?? '',
      status: ProposalStatus.parse(json['status']),
      silentAt: DateTime.parse(json['silentAt']! as String),
      votes: [
        for (final item in json['votes'] as List<Object?>? ?? const [])
          TeamVote.fromJson(item! as Map<String, Object?>),
      ],
    );
  }

  final String id;
  final CoursePerson? proposedBy;
  final TeamWork work;
  final String note;
  final ProposalStatus status;
  final DateTime silentAt;
  final List<TeamVote> votes;

  VoteDecision? decisionOf(String userId) {
    for (final vote in votes) {
      if (vote.person.userId == userId) {
        return vote.decision;
      }
    }
    return null;
  }
}
