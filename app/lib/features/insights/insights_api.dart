import '../../core/api/api_client.dart';
import 'course_topics.dart';

class InsightsApi {
  const InsightsApi(this._api);

  final ApiClient _api;

  Future<CourseTopics> topics(String courseId) async => CourseTopics.fromJson(
    (await _api.get(
          '/v1/navigator/courses/${Uri.encodeComponent(courseId)}/topics',
        ))!
        as Map<String, Object?>,
  );
}
