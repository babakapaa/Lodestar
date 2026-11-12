import '../../core/api/api_client.dart';
import 'account.dart';
import 'device_session.dart';
import 'profile.dart';

class ProfileRepository {
  const ProfileRepository(this._api);

  final ApiClient _api;

  Future<Account> account() async =>
      Account.fromJson((await _api.get('/v1/me'))! as Map<String, Object?>);

  Future<Profile> profile() async => Profile.fromJson(
    (await _api.get('/v1/me/profile'))! as Map<String, Object?>,
  );

  Future<Profile> save(Profile profile) async => Profile.fromJson(
    (await _api.put('/v1/me/profile', profile.toJson()))!
        as Map<String, Object?>,
  );

  Future<List<DeviceSession>> sessions() async => [
    for (final item in (await _api.get('/v1/sessions'))! as List<Object?>)
      DeviceSession.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> revoke(String id) =>
      _api.delete('/v1/sessions/${Uri.encodeComponent(id)}');
}
