import 'dart:typed_data';

import '../../core/api/api_client.dart';
import 'student_import.dart';
import 'table_format.dart';

class GradebookApi {
  const GradebookApi(this._api);

  static const _rejected = 422;

  final ApiClient _api;

  Future<Uint8List> export(String courseId, TableFormat format) =>
      _api.download('${_course(courseId)}/gradebook', {'format': format.name});

  Future<StudentImport> importStudents(
    String courseId,
    Uint8List bytes,
    TableFormat format,
  ) async => StudentImport.fromJson(
    (await _api.upload(
          '${_course(courseId)}/students/import',
          bytes,
          contentType: format.contentType,
          accepted: {_rejected},
        ))!
        as Map<String, Object?>,
  );

  String _course(String id) => '/v1/study/courses/${Uri.encodeComponent(id)}';
}
