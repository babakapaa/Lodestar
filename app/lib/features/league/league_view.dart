import 'award_grant.dart';
import 'division.dart';
import 'league_dispute.dart';
import 'league_entry.dart';
import 'league_json.dart';
import 'league_season.dart';
import 'league_table_row.dart';
import 'league_task.dart';
import 'league_transition.dart';
import 'player_row.dart';
import 'standing.dart';
import 'work_weight.dart';

class LeagueView {
  const LeagueView({
    required this.id,
    required this.kind,
    required this.originId,
    required this.title,
    required this.leads,
    required this.season,
    required this.divisions,
    required this.confirm,
    required this.formula,
    required this.me,
    required this.table,
    required this.entries,
    required this.transitions,
    required this.disputes,
    required this.players,
    required this.pending,
    required this.openDisputes,
    required this.tasks,
    required this.grants,
  });

  factory LeagueView.fromJson(Map<String, Object?> json) {
    final me = json['me'];
    return LeagueView(
      id: json['id']! as String,
      kind: json['kind']! as String,
      originId: json['originId']! as String,
      title: json['title'] as String? ?? '',
      leads: json['role'] == 'lead',
      season: LeagueSeason.fromJson(json['season']! as Map<String, Object?>),
      divisions: LeagueJson.list(json['divisions'], Division.fromJson),
      confirm: json['confirm'] as bool? ?? false,
      formula: LeagueJson.list(json['formula'], WorkWeight.fromJson),
      me: me is Map<String, Object?> ? Standing.fromJson(me) : null,
      table: LeagueJson.list(json['table'], LeagueTableRow.fromJson),
      entries: LeagueJson.list(json['entries'], LeagueEntry.fromJson),
      transitions: LeagueJson.list(
        json['transitions'],
        LeagueTransition.fromJson,
      ),
      disputes: LeagueJson.list(json['disputes'], LeagueDispute.fromJson),
      players: LeagueJson.list(json['players'], PlayerRow.fromJson),
      pending: LeagueJson.list(json['pending'], LeagueTransition.fromJson),
      openDisputes: LeagueJson.list(
        json['openDisputes'],
        LeagueDispute.fromJson,
      ),
      tasks: LeagueJson.list(json['tasks'], LeagueTask.fromJson),
      grants: LeagueJson.list(json['grants'], AwardGrant.fromJson),
    );
  }

  final String id;
  final String kind;
  final String originId;
  final String title;
  final bool leads;
  final LeagueSeason season;
  final List<Division> divisions;
  final bool confirm;
  final List<WorkWeight> formula;
  final Standing? me;
  final List<LeagueTableRow> table;
  final List<LeagueEntry> entries;
  final List<LeagueTransition> transitions;
  final List<LeagueDispute> disputes;
  final List<PlayerRow> players;
  final List<LeagueTransition> pending;
  final List<LeagueDispute> openDisputes;
  final List<LeagueTask> tasks;
  final List<AwardGrant> grants;

  String division(int level) =>
      level >= 0 && level < divisions.length ? divisions[level].name : '';
}
