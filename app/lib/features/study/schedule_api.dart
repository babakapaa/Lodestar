import '../../core/api/api_client.dart';
import 'lesson.dart';
import 'lesson_change.dart';
import 'lesson_series.dart';

class ScheduleApi {
  const ScheduleApi(this._api);

  final ApiClient _api;

  Future<List<LessonSeries>> series(String courseId) async => [
    for (final item
        in (await _api.get('${_course(courseId)}/series'))! as List<Object?>)
      LessonSeries.fromJson(item! as Map<String, Object?>),
  ];

  Future<List<Lesson>> lessons(String courseId, String from, String to) async =>
      [
        for (final item
            in (await _api.get('${_course(courseId)}/schedule', {
                  'from': from,
                  'to': to,
                }))!
                as List<Object?>)
          Lesson.fromJson(item! as Map<String, Object?>),
      ];

  Future<void> create(String courseId, LessonSeries series) =>
      _api.post('${_course(courseId)}/series', series.toJson());

  Future<void> update(LessonSeries series) =>
      _api.put(_series(series.id), series.toJson());

  Future<void> delete(String id) => _api.delete(_series(id));

  Future<void> change(String seriesId, String date, LessonChange change) =>
      _api.put('${_series(seriesId)}/changes/$date', change.toJson());

  Future<void> revert(String seriesId, String date) =>
      _api.delete('${_series(seriesId)}/changes/$date');

  static String _course(String id) =>
      '/v1/study/courses/${Uri.encodeComponent(id)}';

  static String _series(String id) =>
      '/v1/study/series/${Uri.encodeComponent(id)}';
}
