import 'dart:typed_data';

import '../../core/api/api_client.dart';
import 'hardware_report.dart';
import 'pilot_course.dart';
import 'setup_overview.dart';
import 'setup_status.dart';

class SetupApi {
  const SetupApi(this._api);

  final ApiClient _api;

  Future<SetupStatus> status() async => SetupStatus.fromJson(
    (await _api.get('/v1/setup/status'))! as Map<String, Object?>,
  );

  Future<SetupOverview> overview() async => SetupOverview.fromJson(
    (await _api.get('/v1/setup'))! as Map<String, Object?>,
  );

  Future<void> assign(String section, String? userId) => _api.put(
    '/v1/setup/sections/${Uri.encodeComponent(section)}/assignee',
    {'userId': userId},
  );

  Future<void> setSetting(
    String section,
    String key,
    Object? value,
  ) => _api.put(
    '/v1/setup/sections/${Uri.encodeComponent(section)}/settings/${Uri.encodeComponent(key)}',
    {'value': value},
  );

  Future<void> submit(String section) =>
      _api.post('/v1/setup/sections/${Uri.encodeComponent(section)}/submit');

  Future<void> review(String section, {required bool approve, String? note}) =>
      _api.post('/v1/setup/sections/${Uri.encodeComponent(section)}/review', {
        'approve': approve,
        'note': ?note,
      });

  Future<void> answer(String question, String answer, String basis) => _api.put(
    '/v1/setup/legal/${Uri.encodeComponent(question)}',
    {'answer': answer, 'basis': basis},
  );

  Future<void> approve(String question) =>
      _api.post('/v1/setup/legal/${Uri.encodeComponent(question)}/approve');

  Future<Uint8List> protocol(String format) =>
      _api.download('/v1/setup/protocol', {'format': format});

  Future<void> sign() => _api.post('/v1/setup/protocol/sign');

  Future<void> launch(List<String> pilot) =>
      _api.post('/v1/setup/launch', {'pilot': pilot});

  Future<HardwareReport> hardware() async => HardwareReport.fromJson(
    (await _api.get('/v1/setup/hardware'))! as Map<String, Object?>,
  );

  Future<HardwareReport> checkHardware() async => HardwareReport.fromJson(
    (await _api.post('/v1/setup/hardware/check'))! as Map<String, Object?>,
  );

  Future<void> expand() => _api.post('/v1/setup/expand');

  Future<List<PilotCourse>> pilotCourses() async => [
    for (final item
        in (await _api.get('/v1/setup/pilot/courses'))! as List<Object?>)
      PilotCourse.fromJson(item! as Map<String, Object?>),
  ];
}
