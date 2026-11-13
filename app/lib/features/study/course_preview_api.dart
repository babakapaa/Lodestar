import '../../core/api/api_client.dart';
import 'assignment_item.dart';
import 'course_person.dart';

class CoursePreviewApi {
  const CoursePreviewApi(this._api);

  final ApiClient _api;

  Future<List<CoursePerson>> students(String courseId) async => [
    for (final item
        in (await _api.get('${_course(courseId)}/students'))! as List<Object?>)
      CoursePerson.fromJson(item! as Map<String, Object?>),
  ];

  Future<List<AssignmentItem>> preview(
    String courseId, {
    required DateTime at,
    String? studentId,
  }) async => [
    for (final item
        in (await _api.get('${_course(courseId)}/preview', {
              'at': at.toUtc().toIso8601String(),
              'student': ?studentId,
            }))!
            as List<Object?>)
      AssignmentItem.fromJson(item! as Map<String, Object?>),
  ];

  static String _course(String id) =>
      '/v1/study/courses/${Uri.encodeComponent(id)}';
}
