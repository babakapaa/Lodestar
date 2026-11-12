import '../../core/api/api_client.dart';
import 'role.dart';
import 'role_assignment.dart';
import 'role_input.dart';

class RolesRepository {
  const RolesRepository(this._api);

  final ApiClient _api;

  Future<List<Role>> roles() async => [
    for (final item in (await _api.get('/v1/roles'))! as List<Object?>)
      Role.fromJson(item! as Map<String, Object?>),
  ];

  Future<List<String>> catalog() async => [
    for (final item in (await _api.get('/v1/permissions'))! as List<Object?>)
      item! as String,
  ];

  Future<Role> create(RoleInput input) async => Role.fromJson(
    (await _api.post('/v1/roles', input.toJson()))! as Map<String, Object?>,
  );

  Future<Role> update(String id, RoleInput input) async => Role.fromJson(
    (await _api.put('/v1/roles/${_id(id)}', input.toJson()))!
        as Map<String, Object?>,
  );

  Future<void> delete(String id) => _api.delete('/v1/roles/${_id(id)}');

  Future<void> move(String id, int position) =>
      _api.post('/v1/roles/${_id(id)}/move', {'position': position});

  Future<List<RoleAssignment>> assignmentsOf(String userId) async => [
    for (final item
        in (await _api.get('/v1/users/${_id(userId)}/roles'))! as List<Object?>)
      RoleAssignment.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> assign(String roleId, String userId) =>
      _api.post('/v1/role-assignments', {'roleId': roleId, 'userId': userId});

  Future<void> unassign(String assignmentId) =>
      _api.delete('/v1/role-assignments/${_id(assignmentId)}');

  static String _id(String id) => Uri.encodeComponent(id);
}
