import '../study/course_person.dart';
import 'award_grant.dart';
import 'league_dispute.dart';
import 'league_entry.dart';
import 'league_json.dart';
import 'league_transition.dart';
import 'standing.dart';

class PlayerDetails {
  const PlayerDetails({
    required this.person,
    required this.standing,
    required this.entries,
    required this.transitions,
    required this.disputes,
    required this.grants,
  });

  factory PlayerDetails.fromJson(Map<String, Object?> json) => PlayerDetails(
    person: CoursePerson.fromJson(json['person']! as Map<String, Object?>),
    standing: Standing.fromJson(json['standing']! as Map<String, Object?>),
    entries: LeagueJson.list(json['entries'], LeagueEntry.fromJson),
    transitions: LeagueJson.list(
      json['transitions'],
      LeagueTransition.fromJson,
    ),
    disputes: LeagueJson.list(json['disputes'], LeagueDispute.fromJson),
    grants: LeagueJson.list(json['grants'], AwardGrant.fromJson),
  );

  final CoursePerson person;
  final Standing standing;
  final List<LeagueEntry> entries;
  final List<LeagueTransition> transitions;
  final List<LeagueDispute> disputes;
  final List<AwardGrant> grants;
}
