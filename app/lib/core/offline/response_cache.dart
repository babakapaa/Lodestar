import 'dart:convert';

import 'package:shared_preferences/shared_preferences.dart';

class ResponseCache {
  const ResponseCache(this._prefs);

  static const _prefix = 'cache.';

  final SharedPreferences _prefs;

  Object? read(String key) {
    final raw = _prefs.getString('$_prefix$key');
    return raw == null ? null : jsonDecode(raw);
  }

  bool has(String key) => _prefs.containsKey('$_prefix$key');

  Future<void> write(String key, Object? body) =>
      _prefs.setString('$_prefix$key', jsonEncode(body));

  Future<void> clear() async {
    for (final key in _prefs.getKeys()) {
      if (key.startsWith(_prefix)) {
        await _prefs.remove(key);
      }
    }
  }
}
