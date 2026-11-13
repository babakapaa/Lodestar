import '../../core/api/api_client.dart';
import 'award.dart';
import 'award_grant.dart';
import 'award_input.dart';
import 'division.dart';
import 'league_dispute.dart';
import 'league_display.dart';
import 'league_json.dart';
import 'league_solution.dart';
import 'league_task.dart';
import 'league_view.dart';
import 'player_details.dart';
import 'task_details.dart';
import 'task_input.dart';

class LeagueApi {
  const LeagueApi(this._api);

  final ApiClient _api;

  Future<LeagueView> view(String kind, String id) async => _view(
    await _api.get(
      '/v1/league/arenas/${Uri.encodeComponent(kind)}/${Uri.encodeComponent(id)}',
    ),
  );

  Future<LeagueView> rules(
    String leagueId, {
    required List<Division> divisions,
    required bool confirm,
    required List<Map<String, Object?>> weights,
  }) async => _view(
    await _api.put('${_league(leagueId)}/rules', {
      'divisions': [for (final d in divisions) d.toJson()],
      'confirm': confirm,
      'weights': weights,
    }),
  );

  Future<LeagueView> display(
    String leagueId,
    LeagueDisplay display,
    String nickname,
  ) async => _view(
    await _api.put('${_league(leagueId)}/display', {
      'display': display.name,
      'nickname': nickname,
    }),
  );

  Future<LeagueView> closeSeason(String leagueId) async =>
      _view(await _api.post('${_league(leagueId)}/seasons'));

  Future<PlayerDetails> player(String leagueId, String userId) async => _player(
    await _api.get(
      '${_league(leagueId)}/players/${Uri.encodeComponent(userId)}',
    ),
  );

  Future<PlayerDetails> adjust(
    String leagueId,
    String userId,
    int points,
    String reason,
  ) async => _player(
    await _api.post('${_league(leagueId)}/adjustments', {
      'userId': userId,
      'points': points,
      'reason': reason,
    }),
  );

  Future<LeagueDispute> dispute(
    String leagueId,
    String text, {
    String? entryId,
    String? transitionId,
  }) async => LeagueDispute.fromJson(
    (await _api.post('${_league(leagueId)}/disputes', {
          'entryId': entryId,
          'transitionId': transitionId,
          'text': text,
        }))!
        as Map<String, Object?>,
  );

  Future<void> decide(String transitionId, {required bool approve}) =>
      _api.post('/v1/league/transitions/${Uri.encodeComponent(transitionId)}', {
        'approve': approve,
      });

  Future<void> answer(
    String disputeId, {
    required bool accept,
    required String answer,
    int points = 0,
  }) => _api.post('/v1/league/disputes/${Uri.encodeComponent(disputeId)}', {
    'accept': accept,
    'answer': answer,
    'points': points,
  });

  Future<LeagueTask> createTask(String leagueId, TaskInput input) async =>
      LeagueTask.fromJson(
        (await _api.post('${_league(leagueId)}/tasks', input.toJson()))!
            as Map<String, Object?>,
      );

  Future<TaskDetails> task(String id) async => TaskDetails.fromJson(
    (await _api.get(_task(id)))! as Map<String, Object?>,
  );

  Future<LeagueTask> updateTask(String id, TaskInput input) async =>
      LeagueTask.fromJson(
        (await _api.put(_task(id), input.toJson()))! as Map<String, Object?>,
      );

  Future<LeagueSolution> review(
    String solutionId, {
    required int points,
    required List<String> passed,
    required String comment,
  }) async => LeagueSolution.fromJson(
    (await _api.post('${_solution(solutionId)}/review', {
          'points': points,
          'passed': passed,
          'comment': comment,
        }))!
        as Map<String, Object?>,
  );

  Future<LeagueSolution> explain(String solutionId, String text) async =>
      LeagueSolution.fromJson(
        (await _api.put('${_solution(solutionId)}/explanation', {
              'text': text,
            }))!
            as Map<String, Object?>,
      );

  Future<List<Award>> awards() async =>
      LeagueJson.list(await _api.get('/v1/league/awards'), Award.fromJson);

  Future<Award> createAward(AwardInput input) async => Award.fromJson(
    (await _api.post('/v1/league/awards', input.toJson()))!
        as Map<String, Object?>,
  );

  Future<Award> updateAward(String id, AwardInput input) async =>
      Award.fromJson(
        (await _api.put(
              '/v1/league/awards/${Uri.encodeComponent(id)}',
              input.toJson(),
            ))!
            as Map<String, Object?>,
      );

  Future<AwardGrant> grant(
    String leagueId,
    String awardId,
    String userId,
    String reason,
  ) async => AwardGrant.fromJson(
    (await _api.post('${_league(leagueId)}/grants', {
          'awardId': awardId,
          'userId': userId,
          'reason': reason,
        }))!
        as Map<String, Object?>,
  );

  Future<void> revoke(String grantId) =>
      _api.delete('/v1/league/grants/${Uri.encodeComponent(grantId)}');

  static LeagueView _view(Object? body) =>
      LeagueView.fromJson(body! as Map<String, Object?>);

  static PlayerDetails _player(Object? body) =>
      PlayerDetails.fromJson(body! as Map<String, Object?>);

  static String _league(String id) =>
      '/v1/league/leagues/${Uri.encodeComponent(id)}';

  static String _task(String id) =>
      '/v1/league/tasks/${Uri.encodeComponent(id)}';

  static String _solution(String id) =>
      '/v1/league/solutions/${Uri.encodeComponent(id)}';
}
