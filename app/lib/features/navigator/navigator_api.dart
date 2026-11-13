import '../../core/api/api_client.dart';
import '../../core/offline/cached_reader.dart';
import '../../core/offline/loaded.dart';
import 'course_material.dart';
import 'dialog_details.dart';
import 'dialog_scope.dart';
import 'material_input.dart';
import 'navigator_answer.dart';
import 'navigator_dialog.dart';
import 'question.dart';

class NavigatorApi {
  const NavigatorApi(this._api, this._reader);

  final ApiClient _api;
  final CachedReader _reader;

  Future<Loaded<List<NavigatorDialog>>> dialogs() async {
    final cached = await _reader.get('/v1/navigator/dialogs');
    return Loaded([
      for (final item in cached.body! as List<Object?>)
        NavigatorDialog.fromJson(item! as Map<String, Object?>),
    ], stale: cached.stale);
  }

  Future<Loaded<DialogDetails>> dialog(String id) async {
    final cached = await _reader.get(_dialog(id));
    return Loaded(
      DialogDetails.fromJson(cached.body! as Map<String, Object?>),
      stale: cached.stale,
    );
  }

  Future<DialogDetails> create(String id, DialogScope? scope) async =>
      DialogDetails.fromJson(
        (await _api.post('/v1/navigator/dialogs', {
              'id': id,
              ...?scope?.toJson(),
            }))!
            as Map<String, Object?>,
      );

  Future<NavigatorAnswer> ask(String dialogId, Question question) async =>
      NavigatorAnswer.fromJson(
        (await _api.post('${_dialog(dialogId)}/messages', question.toJson()))!
            as Map<String, Object?>,
      );

  Future<void> delete(String id) => _api.delete(_dialog(id));

  Future<Loaded<List<CourseMaterial>>> materials(String courseId) async {
    final cached = await _reader.get(_materials(courseId));
    return Loaded([
      for (final item in cached.body! as List<Object?>)
        CourseMaterial.fromJson(item! as Map<String, Object?>),
    ], stale: cached.stale);
  }

  Future<void> addMaterial(String courseId, MaterialInput input) =>
      _api.post(_materials(courseId), input.toJson());

  Future<void> deleteMaterial(String id) =>
      _api.delete('/v1/navigator/materials/${Uri.encodeComponent(id)}');

  static String _dialog(String id) =>
      '/v1/navigator/dialogs/${Uri.encodeComponent(id)}';

  static String _materials(String courseId) =>
      '/v1/navigator/courses/${Uri.encodeComponent(courseId)}/materials';
}
