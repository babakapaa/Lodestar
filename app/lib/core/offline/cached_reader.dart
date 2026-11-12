import '../api/api_client.dart';
import '../api/api_exception.dart';
import 'cached_body.dart';
import 'response_cache.dart';

class CachedReader {
  const CachedReader(this._api, this._cache);

  final ApiClient _api;
  final ResponseCache _cache;

  Future<CachedBody> get(String path, [Map<String, String>? query]) async {
    final key = Uri(path: path, queryParameters: query).toString();
    try {
      final body = await _api.get(path, query);
      await _cache.write(key, body);
      return CachedBody(body, stale: false);
    } on ApiException catch (error) {
      if (error.code != ApiException.unreachable || !_cache.has(key)) {
        rethrow;
      }
      return CachedBody(_cache.read(key), stale: true);
    }
  }
}
