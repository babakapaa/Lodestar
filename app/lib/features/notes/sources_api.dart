import '../../core/api/api_client.dart';
import '../../core/offline/cached_reader.dart';
import '../../core/offline/client_id.dart';
import '../../core/offline/loaded.dart';
import 'source.dart';

class SourcesApi {
  const SourcesApi(this._api, this._reader);

  final ApiClient _api;
  final CachedReader _reader;

  Future<Loaded<List<Source>>> list() async {
    final cached = await _reader.get('/v1/notes/sources');
    return Loaded([
      for (final item in cached.body! as List<Object?>)
        Source.fromJson(item! as Map<String, Object?>),
    ], stale: cached.stale);
  }

  Future<void> save(Source source) =>
      _api.post('/v1/notes/sources', source.toJson());

  Future<void> delete(String id) =>
      _api.delete('/v1/notes/sources/${Uri.encodeComponent(id)}');

  Future<Source> lookup({String? doi, String? url}) async => Source.fromJson(
    (await _api.post('/v1/notes/lookup', {'doi': ?doi, 'url': ?url}))!
        as Map<String, Object?>,
    id: ClientId.generate(),
  );
}
