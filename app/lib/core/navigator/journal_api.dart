import '../api/api_client.dart';
import 'journal_entry.dart';

class JournalApi {
  const JournalApi(this._api);

  static const preference = 'navigator.journal';

  final ApiClient _api;

  Future<List<JournalEntry>> list({String? scopeKind, String? scopeId}) async {
    final body = await _api.get('/v1/navigator/journal', {
      'scopeKind': ?scopeKind,
      'scopeId': ?scopeId,
    });
    return [
      for (final item in body! as List<Object?>)
        JournalEntry.fromJson(item! as Map<String, Object?>),
    ];
  }

  Future<void> delete(String id) =>
      _api.delete('/v1/navigator/journal/${Uri.encodeComponent(id)}');
}
