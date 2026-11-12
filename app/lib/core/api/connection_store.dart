import 'package:shared_preferences/shared_preferences.dart';

import 'server_address.dart';

class ConnectionStore {
  const ConnectionStore(this._prefs);

  static const _serverKey = 'connection.server';
  static const _tenantKey = 'connection.tenant';
  static const presetServer = String.fromEnvironment('LODESTAR_SERVER');
  static const presetTenant = String.fromEnvironment('LODESTAR_TENANT');

  final SharedPreferences _prefs;

  ServerAddress? load() {
    final server = _prefs.getString(_serverKey);
    final tenant = _prefs.getString(_tenantKey);
    if (server != null && tenant != null) {
      return ServerAddress.tryParse(server, tenant);
    }
    return ServerAddress.tryParse(presetServer, presetTenant);
  }

  Future<void> save(ServerAddress address) async {
    await Future.wait([
      _prefs.setString(_serverKey, address.server.toString()),
      _prefs.setString(_tenantKey, address.tenant),
    ]);
  }

  Future<void> clear() async {
    await Future.wait([_prefs.remove(_serverKey), _prefs.remove(_tenantKey)]);
  }
}
