import '../study/course_person.dart';
import 'history_entry.dart';
import 'leave_request.dart';
import 'team.dart';
import 'team_people.dart';
import 'team_proposal.dart';
import 'team_scope.dart';

class TeamDetails {
  const TeamDetails({
    required this.team,
    required this.scope,
    required this.manage,
    required this.inviteCode,
    required this.invited,
    required this.history,
    required this.request,
    required this.canLeaveFreely,
    required this.proposal,
  });

  factory TeamDetails.fromJson(Map<String, Object?> json) {
    final request = json['request'];
    final proposal = json['proposal'];
    return TeamDetails(
      team: Team.fromJson(json),
      scope: TeamScope.fromJson(json['scope']! as Map<String, Object?>),
      manage: json['manage'] as bool? ?? false,
      inviteCode: json['inviteCode'] as String?,
      invited: TeamPeople.parse(json['invited']),
      history: [
        for (final item in json['history'] as List<Object?>? ?? const [])
          HistoryEntry.fromJson(item! as Map<String, Object?>),
      ],
      request: request is Map<String, Object?>
          ? LeaveRequest.fromJson(request)
          : null,
      canLeaveFreely: json['canLeaveFreely'] as bool? ?? false,
      proposal: proposal is Map<String, Object?>
          ? TeamProposal.fromJson(proposal)
          : null,
    );
  }

  final Team team;
  final TeamScope scope;
  final bool manage;
  final String? inviteCode;
  final List<CoursePerson> invited;
  final List<HistoryEntry> history;
  final LeaveRequest? request;
  final bool canLeaveFreely;
  final TeamProposal? proposal;
}
