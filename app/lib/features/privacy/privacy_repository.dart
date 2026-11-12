import 'dart:typed_data';

import '../../core/api/api_client.dart';
import '../../core/api/api_exception.dart';
import 'deletion.dart';

class PrivacyRepository {
  const PrivacyRepository(this._api);

  final ApiClient _api;

  Future<Uint8List> export() => _api.download('/v1/me/export');

  Future<Deletion?> deletion() async {
    try {
      final body = await _api.get('/v1/me/deletion');
      return Deletion.fromJson(body! as Map<String, Object?>);
    } on ApiException catch (error) {
      if (error.statusCode == 404) {
        return null;
      }
      rethrow;
    }
  }

  Future<Deletion> requestDeletion() async {
    final body = await _api.post('/v1/me/deletion');
    return Deletion.fromJson(body! as Map<String, Object?>);
  }

  Future<void> cancelDeletion() => _api.delete('/v1/me/deletion');
}
