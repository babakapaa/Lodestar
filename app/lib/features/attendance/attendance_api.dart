import 'dart:typed_data';

import '../../core/api/api_client.dart';
import '../study/table_format.dart';
import 'attendance_check.dart';
import 'attendance_json.dart';
import 'attendance_qr.dart';
import 'attendance_session.dart';
import 'attendance_status.dart';
import 'course_journal.dart';
import 'open_lesson.dart';

class AttendanceApi {
  const AttendanceApi(this._api);

  final ApiClient _api;

  static const _base = '/v1/attendance';

  Future<AttendanceSession> open(String seriesId, String date) async =>
      _session(
        await _api.post(
          '$_base/lessons/${Uri.encodeComponent(seriesId)}/$date/open',
        ),
      );

  Future<AttendanceSession> session(String id) async =>
      _session(await _api.get(_path(id)));

  Future<AttendanceQr> qr(String id) async => AttendanceQr.fromJson(
    (await _api.get('${_path(id)}/qr'))! as Map<String, Object?>,
  );

  Future<AttendanceSession> refreshCode(String id) async =>
      _session(await _api.post('${_path(id)}/code'));

  Future<AttendanceSession> mark(
    String id,
    String userId,
    AttendanceStatus? status,
  ) async => _session(
    await _api.put('${_path(id)}/marks/${Uri.encodeComponent(userId)}', {
      'status': status?.name,
    }),
  );

  Future<AttendanceSession> close(String id) async =>
      _session(await _api.post('${_path(id)}/close'));

  Future<AttendanceCheck> check({
    String? sessionId,
    String? token,
    String? code,
  }) async => AttendanceCheck.fromJson(
    (await _api.post('$_base/check', {
          'sessionId': ?sessionId,
          'token': ?token,
          'code': ?code,
        }))!
        as Map<String, Object?>,
  );

  Future<List<OpenLesson>> openLessons() async =>
      AttendanceJson.list(await _api.get('$_base/open'), OpenLesson.fromJson);

  Future<CourseJournal> journal(String courseId) async =>
      CourseJournal.fromJson(
        (await _api.get('$_base/courses/${Uri.encodeComponent(courseId)}'))!
            as Map<String, Object?>,
      );

  Future<Uint8List> export(String courseId, TableFormat format) =>
      _api.download('$_base/courses/${Uri.encodeComponent(courseId)}/export', {
        'format': format.name,
      });

  String _path(String id) => '$_base/sessions/${Uri.encodeComponent(id)}';

  AttendanceSession _session(Object? json) =>
      AttendanceSession.fromJson(json! as Map<String, Object?>);
}
