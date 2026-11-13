import '../../core/api/api_client.dart';
import 'leave_result.dart';
import 'team_card.dart';
import 'team_details.dart';
import 'team_profile.dart';
import 'team_scope_view.dart';
import 'team_suggestion.dart';

class TeamApi {
  const TeamApi(this._api);

  final ApiClient _api;

  Future<List<TeamCard>> mine() => _cards('/v1/team/teams');

  Future<List<TeamCard>> invitations() => _cards('/v1/team/invitations');

  Future<TeamScopeView> scope(String kind, String id) async =>
      _view(await _api.get(_scope(kind, id)));

  Future<TeamDetails> create(
    String kind,
    String id,
    String name, [
    List<String> members = const [],
  ]) async => _details(
    await _api.post('${_scope(kind, id)}/teams', {
      'name': name,
      'members': members,
    }),
  );

  Future<TeamScopeView> seek(
    String kind,
    String id, {
    required bool seeking,
    String note = '',
  }) async => _view(
    await _api.put('${_scope(kind, id)}/seeking', {
      'seeking': seeking,
      'note': note,
    }),
  );

  Future<TeamScopeView> distribute(String kind, String id) async =>
      _view(await _api.post('${_scope(kind, id)}/distribution'));

  Future<TeamSuggestion> suggest(String kind, String id) async =>
      TeamSuggestion.fromJson(
        (await _api.post('${_scope(kind, id)}/suggestions'))!
            as Map<String, Object?>,
      );

  Future<TeamDetails> get(String id) async =>
      _details(await _api.get(_team(id)));

  Future<TeamDetails> rename(String id, String name) async =>
      _details(await _api.patch(_team(id), {'name': name}));

  Future<TeamDetails> label(String id, String label) async =>
      _details(await _api.put('${_team(id)}/label', {'label': label}));

  Future<TeamDetails> invite(String id, String userId) async =>
      _details(await _api.post('${_team(id)}/invitations', {'userId': userId}));

  Future<void> uninvite(String id, String userId) =>
      _api.delete('${_team(id)}/invitations/${Uri.encodeComponent(userId)}');

  Future<TeamDetails> accept(String id) async =>
      _details(await _api.post('${_team(id)}/accept'));

  Future<LeaveResult> leave(String id, String reason) async =>
      LeaveResult.fromJson(
        (await _api.post('${_team(id)}/leave', {'reason': reason}))!
            as Map<String, Object?>,
      );

  Future<TeamDetails> add(String id, String userId) async => _details(
    await _api.put('${_team(id)}/members/${Uri.encodeComponent(userId)}', null),
  );

  Future<TeamDetails> remove(
    String id,
    String userId, {
    bool solo = false,
  }) async => _details(
    await _api.delete(
      '${_team(id)}/members/${Uri.encodeComponent(userId)}',
      solo ? {'solo': 'true'} : null,
    ),
  );

  Future<TeamCard> preview(String code) async =>
      TeamCard.fromJson((await _api.get(_code(code)))! as Map<String, Object?>);

  Future<TeamDetails> join(String code) async =>
      _details(await _api.post(_code(code)));

  Future<void> decide(String requestId, {required bool approve}) => _api.post(
    '/v1/team/leave-requests/${Uri.encodeComponent(requestId)}',
    {'approve': approve},
  );

  Future<TeamDetails> propose(
    String id,
    Map<String, Object?> work,
    String note,
  ) async => _details(
    await _api.post('${_team(id)}/proposals', {'work': work, 'note': note}),
  );

  Future<TeamDetails> vote(
    String proposalId, {
    required bool agree,
    String comment = '',
  }) async => _details(
    await _api.post(
      '/v1/team/proposals/${Uri.encodeComponent(proposalId)}/votes',
      {'agree': agree, 'comment': comment},
    ),
  );

  Future<TeamProfile> profile() async => TeamProfile.fromJson(
    (await _api.get('/v1/team/card'))! as Map<String, Object?>,
  );

  Future<TeamProfile> saveProfile(TeamProfile profile) async =>
      TeamProfile.fromJson(
        (await _api.put('/v1/team/card', profile.toJson()))!
            as Map<String, Object?>,
      );

  Future<List<TeamCard>> _cards(String path) async => [
    for (final item in (await _api.get(path))! as List<Object?>)
      TeamCard.fromJson(item! as Map<String, Object?>),
  ];

  static TeamScopeView _view(Object? body) =>
      TeamScopeView.fromJson(body! as Map<String, Object?>);

  static TeamDetails _details(Object? body) =>
      TeamDetails.fromJson(body! as Map<String, Object?>);

  static String _scope(String kind, String id) =>
      '/v1/team/scopes/${Uri.encodeComponent(kind)}/${Uri.encodeComponent(id)}';

  static String _team(String id) => '/v1/team/teams/${Uri.encodeComponent(id)}';

  static String _code(String code) =>
      '/v1/team/codes/${Uri.encodeComponent(code)}';
}
