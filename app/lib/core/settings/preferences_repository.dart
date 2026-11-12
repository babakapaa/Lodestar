import '../api/api_client.dart';
import 'preference_value.dart';

class PreferencesRepository {
  const PreferencesRepository(this._api);

  final ApiClient _api;

  Future<Map<String, PreferenceValue>> values() async {
    final body = await _api.get('/v1/me/preferences');
    return {
      for (final item in body! as List<Object?>)
        if (item case {
          'key': final String key,
          'value': final Object? value,
          'set': final bool set,
        })
          key: PreferenceValue(value: value, set: set),
    };
  }

  Future<void> set(String key, Object? value) => _api.put(
    '/v1/me/preferences/${Uri.encodeComponent(key)}',
    {'value': value},
  );
}
