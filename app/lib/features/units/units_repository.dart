import '../../core/api/api_client.dart';
import 'unit.dart';
import 'unit_details.dart';
import 'unit_input.dart';
import 'unit_member.dart';

class UnitsRepository {
  const UnitsRepository(this._api);

  final ApiClient _api;

  Future<List<Unit>> children(
    String? parentId, {
    required bool archived,
  }) async {
    final body = await _api.get('/v1/units', {
      'parent': ?parentId,
      if (archived) 'archived': 'true',
    });
    return [
      for (final item in body! as List<Object?>)
        Unit.fromJson(item! as Map<String, Object?>),
    ];
  }

  Future<UnitDetails> details(String id) async => UnitDetails.fromJson(
    (await _api.get('/v1/units/${_id(id)}'))! as Map<String, Object?>,
  );

  Future<void> create(UnitInput input) =>
      _api.post('/v1/units', input.toJson());

  Future<void> update(String id, UnitInput input) =>
      _api.put('/v1/units/${_id(id)}', input.toJson());

  Future<void> move(String id, String? parentId) =>
      _api.post('/v1/units/${_id(id)}/move', {'parentId': parentId});

  Future<void> archive(String id) => _api.post('/v1/units/${_id(id)}/archive');

  Future<void> restore(String id) => _api.post('/v1/units/${_id(id)}/restore');

  Future<List<UnitMember>> members(String id) async => [
    for (final item
        in (await _api.get('/v1/units/${_id(id)}/members'))! as List<Object?>)
      UnitMember.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> addMember(String id, String userId) =>
      _api.put('/v1/units/${_id(id)}/members/${_id(userId)}', null);

  Future<void> removeMember(String id, String userId) =>
      _api.delete('/v1/units/${_id(id)}/members/${_id(userId)}');

  static String _id(String id) => Uri.encodeComponent(id);
}
