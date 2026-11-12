import '../../core/api/api_client.dart';
import 'person_hit.dart';

class PeopleRepository {
  const PeopleRepository(this._api);

  final ApiClient _api;

  Future<List<PersonHit>> search(String query) async {
    final text = query.trim();
    if (text.isEmpty) {
      return const [];
    }
    final body = await _api.get('/v1/search', {'q': text, 'kind': 'people'});
    return [
      for (final item in body! as List<Object?>)
        PersonHit.fromJson(item! as Map<String, Object?>),
    ];
  }

  Future<String> name(String userId) async {
    final body = await _api.get(
      '/v1/users/${Uri.encodeComponent(userId)}/profile',
    );
    return (body! as Map<String, Object?>)['displayName']! as String;
  }
}
