import '../../core/api/api_client.dart';
import '../../core/offline/cached_body.dart';
import '../../core/offline/cached_reader.dart';
import '../../core/offline/loaded.dart';
import '../../core/sources/source.dart';
import 'agenda.dart';
import 'assignment.dart';
import 'assignment_details.dart';
import 'assignment_input.dart';
import 'assignment_item.dart';
import 'course.dart';
import 'course_details.dart';
import 'grade_input.dart';
import 'lesson.dart';
import 'search_map.dart';
import 'search_map_draft.dart';
import 'student_work.dart';
import 'submission.dart';
import 'submission_draft.dart';

class StudyRepository {
  const StudyRepository(this._api, this._reader);

  final ApiClient _api;
  final CachedReader _reader;

  Future<Loaded<List<Course>>> courses() async =>
      _list(await _reader.get('/v1/study/courses'), Course.fromJson);

  Future<Loaded<CourseDetails>> course(String id) async => _one(
    await _reader.get('/v1/study/courses/${_id(id)}'),
    CourseDetails.fromJson,
  );

  Future<Loaded<List<AssignmentItem>>> assignments(String courseId) async =>
      _list(
        await _reader.get('/v1/study/courses/${_id(courseId)}/assignments'),
        AssignmentItem.fromJson,
      );

  Future<Loaded<AssignmentDetails>> assignment(String id) async => _one(
    await _reader.get('/v1/study/assignments/${_id(id)}'),
    AssignmentDetails.fromJson,
  );

  Future<Assignment> create(String courseId, AssignmentInput input) async =>
      Assignment.fromJson(
        (await _api.post(
              '/v1/study/courses/${_id(courseId)}/assignments',
              input.toJson(),
            ))!
            as Map<String, Object?>,
      );

  Future<Assignment> update(String id, AssignmentInput input) async =>
      Assignment.fromJson(
        (await _api.put('/v1/study/assignments/${_id(id)}', input.toJson()))!
            as Map<String, Object?>,
      );

  Future<void> delete(String id) =>
      _api.delete('/v1/study/assignments/${_id(id)}');

  Future<void> submit(SubmissionDraft draft) => _api.post(
    '/v1/study/assignments/${_id(draft.assignmentId)}/submissions',
    draft.toRequest(),
  );

  Future<Loaded<SearchMapDraft>> searchMap(String assignmentId) async => _one(
    await _reader.get('/v1/study/assignments/${_id(assignmentId)}/search-map'),
    SearchMapDraft.fromJson,
  );

  Future<SearchMapDraft> saveSearchMap(
    String assignmentId,
    SearchMap map,
  ) async => SearchMapDraft.fromJson(
    (await _api.put(
          '/v1/study/assignments/${_id(assignmentId)}/search-map',
          map.toJson(),
        ))!
        as Map<String, Object?>,
  );

  Future<Loaded<List<Source>>> sources() async =>
      _list(await _reader.get('/v1/notes/sources'), Source.fromJson);

  Future<List<StudentWork>> work(String assignmentId) async => [
    for (final item
        in (await _api.get('/v1/study/assignments/${_id(assignmentId)}/work'))!
            as List<Object?>)
      StudentWork.fromJson(item! as Map<String, Object?>),
  ];

  Future<List<Submission>> attempts(
    String assignmentId,
    String studentId,
  ) async => [
    for (final item
        in (await _api.get(
              '/v1/study/assignments/${_id(assignmentId)}/students/${_id(studentId)}/submissions',
            ))!
            as List<Object?>)
      Submission.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> grade(String assignmentId, String studentId, GradeInput input) =>
      _api.put(
        '/v1/study/assignments/${_id(assignmentId)}/grades/${_id(studentId)}',
        input.toJson(),
      );

  Future<Loaded<Agenda>> agenda(String? date) async => _one(
    await _reader.get('/v1/study/agenda', {'date': ?date}),
    Agenda.fromJson,
  );

  Future<Loaded<List<Lesson>>> schedule(String from, String to) async => _list(
    await _reader.get('/v1/study/schedule', {'from': from, 'to': to}),
    Lesson.fromJson,
  );

  Loaded<List<T>> _list<T>(
    CachedBody cached,
    T Function(Map<String, Object?>) parse,
  ) => Loaded([
    for (final item in cached.body! as List<Object?>)
      parse(item! as Map<String, Object?>),
  ], stale: cached.stale);

  Loaded<T> _one<T>(
    CachedBody cached,
    T Function(Map<String, Object?>) parse,
  ) => Loaded(parse(cached.body! as Map<String, Object?>), stale: cached.stale);

  static String _id(String id) => Uri.encodeComponent(id);
}
