import 'dart:convert';
import 'dart:typed_data';

import '../../core/api/api_client.dart';
import 'note_changes.dart';
import 'note_member.dart';
import 'note_meta.dart';
import 'remote_note.dart';

class NotesApi {
  const NotesApi(this._api);

  final ApiClient _api;

  Future<List<RemoteNote>> list() async => [
    for (final item in (await _api.get('/v1/notes'))! as List<Object?>)
      RemoteNote.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> create(String id, NoteMeta meta) =>
      _api.post('/v1/notes', {'id': id, 'meta': meta.toJson()});

  Future<void> setMeta(String id, NoteMeta meta) =>
      _api.put('/v1/notes/${_id(id)}/meta', meta.toJson());

  Future<void> delete(String id) => _api.delete('/v1/notes/${_id(id)}');

  Future<void> restore(String id) => _api.post('/v1/notes/${_id(id)}/restore');

  Future<NoteChanges> changes(String id, int after) async =>
      NoteChanges.fromJson(
        (await _api.get('/v1/notes/${_id(id)}/updates', {'after': '$after'}))!
            as Map<String, Object?>,
      );

  Future<void> push(String id, String key, Uint8List data) => _api.post(
    '/v1/notes/${_id(id)}/updates',
    {'key': key, 'data': base64Encode(data)},
  );

  Future<void> compact(String id, int seq, Uint8List data, String body) =>
      _api.post('/v1/notes/${_id(id)}/versions', {
        'seq': seq,
        'data': base64Encode(data),
        'body': body,
      });

  Future<List<NoteMember>> members(String id) async => [
    for (final item
        in (await _api.get('/v1/notes/${_id(id)}/members'))! as List<Object?>)
      NoteMember.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> share(String id, String userId, {required bool canEdit}) => _api
      .put('/v1/notes/${_id(id)}/members/${_id(userId)}', {'canEdit': canEdit});

  Future<void> unshare(String id, String userId) =>
      _api.delete('/v1/notes/${_id(id)}/members/${_id(userId)}');

  static String _id(String id) => Uri.encodeComponent(id);
}
