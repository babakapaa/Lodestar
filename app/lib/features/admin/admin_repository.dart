import '../../core/api/api_client.dart';
import '../../core/settings/preference_value.dart';
import 'audit_entry.dart';
import 'module_state.dart';
import 'setting_definition.dart';

class AdminRepository {
  const AdminRepository(this._api);

  final ApiClient _api;

  Future<List<SettingDefinition>> schema() async => [
    for (final item
        in (await _api.get('/v1/settings/schema'))! as List<Object?>)
      SettingDefinition.fromJson(item! as Map<String, Object?>),
  ];

  Future<Map<String, PreferenceValue>> values() async {
    final body = await _api.get('/v1/settings');
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

  Future<void> set(String key, Object? value) =>
      _api.put('/v1/settings/${Uri.encodeComponent(key)}', {'value': value});

  Future<void> reset(String key) =>
      _api.delete('/v1/settings/${Uri.encodeComponent(key)}');

  Future<List<ModuleState>> modules() async => [
    for (final item in (await _api.get('/v1/modules'))! as List<Object?>)
      ModuleState.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> setModule(String name, {required bool enabled}) => _api.put(
    '/v1/modules/${Uri.encodeComponent(name)}',
    {'enabled': enabled},
  );

  Future<List<AuditEntry>> audit({int? before}) async {
    final body = await _api.get('/v1/audit', {
      if (before != null) 'before': '$before',
    });
    return [
      for (final item in body! as List<Object?>)
        AuditEntry.fromJson(item! as Map<String, Object?>),
    ];
  }
}
