import 'dart:typed_data';

import '../../core/api/api_client.dart';
import '../study/table_format.dart';
import 'gpu_class.dart';
import 'gpu_node.dart';
import 'run_estimate.dart';
import 'run_input.dart';
import 'run_kind.dart';
import 'sandbox_json.dart';
import 'sandbox_quota.dart';
import 'sandbox_run.dart';
import 'sandbox_template.dart';
import 'sandbox_term.dart';
import 'topup.dart';

class SandboxApi {
  const SandboxApi(this._api);

  final ApiClient _api;

  static const _base = '/v1/sandbox';

  Future<SandboxQuota> quota() async => SandboxQuota.fromJson(
    (await _api.get('$_base/quota'))! as Map<String, Object?>,
  );

  Future<void> transfer(String teamId, double hours) =>
      _api.post('$_base/transfers', {'teamId': teamId, 'hours': hours});

  Future<List<Topup>> topups() async =>
      SandboxJson.list(await _api.get('$_base/topups'), Topup.fromJson);

  Future<void> requestTopup({
    required String workKind,
    required String workId,
    required double hours,
    required String reason,
  }) => _api.post('$_base/topups', {
    'workKind': workKind,
    'workId': workId,
    'hours': hours,
    'reason': reason,
  });

  Future<void> decideTopup(String id, {required bool approve, double? hours}) =>
      _api.post('$_base/topups/${Uri.encodeComponent(id)}/decision', {
        'approve': approve,
        'hours': ?hours,
      });

  Future<List<SandboxRun>> runs() async =>
      SandboxJson.list(await _api.get('$_base/runs'), SandboxRun.fromJson);

  Future<SandboxRun> run(String id) async => _run(await _api.get(_runPath(id)));

  Future<SandboxRun> start(RunInput input) async =>
      _run(await _api.post('$_base/runs', input.toJson()));

  Future<SandboxRun> stop(String id) async =>
      _run(await _api.post('${_runPath(id)}/stop'));

  Future<SandboxRun> rerun(String id) async =>
      _run(await _api.post('${_runPath(id)}/rerun'));

  Future<RunEstimate> estimate({
    required RunKind kind,
    required String templateId,
    double? maxHours,
    String? teamId,
  }) async => RunEstimate.fromJson(
    (await _api.post('$_base/estimate', {
          'kind': kind.name,
          'templateId': templateId,
          'maxHours': ?maxHours,
          'teamId': ?teamId,
        }))!
        as Map<String, Object?>,
  );

  Future<Uint8List> export(TableFormat format) =>
      _api.download('$_base/runs/export', {'format': format.name});

  Future<List<SandboxRun>> held() async =>
      SandboxJson.list(await _api.get('$_base/runs/held'), SandboxRun.fromJson);

  Future<void> decideHold(String id, {required bool release}) =>
      _api.post('${_runPath(id)}/hold-decision', {'release': release});

  Future<List<GpuClass>> classes() async =>
      SandboxJson.list(await _api.get('$_base/classes'), GpuClass.fromJson);

  Future<void> saveClass(GpuClass value) => value.id.isEmpty
      ? _api.post('$_base/classes', value.toJson())
      : _api.put(
          '$_base/classes/${Uri.encodeComponent(value.id)}',
          value.toJson(),
        );

  Future<List<SandboxTerm>> terms() async =>
      SandboxJson.list(await _api.get('$_base/terms'), SandboxTerm.fromJson);

  Future<void> saveTerm(SandboxTerm value) => value.id.isEmpty
      ? _api.post('$_base/terms', value.toJson())
      : _api.put(
          '$_base/terms/${Uri.encodeComponent(value.id)}',
          value.toJson(),
        );

  Future<void> deleteTerm(String id) =>
      _api.delete('$_base/terms/${Uri.encodeComponent(id)}');

  Future<List<SandboxTemplate>> templates() async => SandboxJson.list(
    await _api.get('$_base/templates'),
    SandboxTemplate.fromJson,
  );

  Future<void> saveTemplate(SandboxTemplate value) => value.id.isEmpty
      ? _api.post('$_base/templates', value.toJson())
      : _api.put(
          '$_base/templates/${Uri.encodeComponent(value.id)}',
          value.toJson(),
        );

  Future<List<GpuNode>> nodes() async =>
      SandboxJson.list(await _api.get('$_base/nodes'), GpuNode.fromJson);

  Future<GpuNode> createNode({
    required String name,
    required String classId,
    required int gpus,
  }) async => GpuNode.fromJson(
    (await _api.post('$_base/nodes', {
          'name': name,
          'classId': classId,
          'gpus': gpus,
        }))!
        as Map<String, Object?>,
  );

  Future<void> deleteNode(String id) =>
      _api.delete('$_base/nodes/${Uri.encodeComponent(id)}');

  String _runPath(String id) => '$_base/runs/${Uri.encodeComponent(id)}';

  SandboxRun _run(Object? json) =>
      SandboxRun.fromJson(json! as Map<String, Object?>);
}
