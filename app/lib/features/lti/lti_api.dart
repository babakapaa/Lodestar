import '../../core/api/api_client.dart';
import 'assignment_source.dart';
import 'link_target.dart';
import 'lti_deep_link.dart';
import 'lti_context.dart';
import 'lti_platform.dart';
import 'lti_platform_input.dart';
import 'lti_tool.dart';
import 'sync_report.dart';

class LtiApi {
  const LtiApi(this._api);

  final ApiClient _api;

  Future<LtiTool> tool() async => LtiTool.fromJson(
    (await _api.get('/v1/lti/tool'))! as Map<String, Object?>,
  );

  Future<List<LtiPlatform>> platforms() async => [
    for (final item in (await _api.get('/v1/lti/platforms'))! as List<Object?>)
      LtiPlatform.fromJson(item! as Map<String, Object?>),
  ];

  Future<LtiPlatform> create(LtiPlatformInput input) async =>
      LtiPlatform.fromJson(
        (await _api.post('/v1/lti/platforms', input.toJson()))!
            as Map<String, Object?>,
      );

  Future<LtiPlatform> update(String id, LtiPlatformInput input) async =>
      LtiPlatform.fromJson(
        (await _api.put(
              '/v1/lti/platforms/${Uri.encodeComponent(id)}',
              input.toJson(),
            ))!
            as Map<String, Object?>,
      );

  Future<void> delete(String id) =>
      _api.delete('/v1/lti/platforms/${Uri.encodeComponent(id)}');

  Future<LtiContext> context(String id) async =>
      _context(await _api.get(_path(id)));

  Future<List<LtiContext>> byOrigin(String kind, String id) async => [
    for (final item
        in (await _api.get(
              '/v1/lti/links/${Uri.encodeComponent(kind)}/${Uri.encodeComponent(id)}',
            ))!
            as List<Object?>)
      LtiContext.fromJson(item! as Map<String, Object?>),
  ];

  Future<LtiContext> link(String id, LinkTarget target) async => _context(
    await _api.put('${_path(id)}/link', {'kind': target.kind, 'id': target.id}),
  );

  Future<LtiContext> unlink(String id) async =>
      _context(await _api.delete('${_path(id)}/link'));

  Future<LtiContext> source(String id, AssignmentSource? source) async =>
      _context(await _api.put('${_path(id)}/source', {'source': source?.name}));

  Future<SyncReport> sync(String id) async => SyncReport.fromJson(
    (await _api.post('${_path(id)}/sync'))! as Map<String, Object?>,
  );

  Future<LtiDeepLink> deepLink(String id) async => LtiDeepLink.fromJson(
    (await _api.get(_deepLinkPath(id)))! as Map<String, Object?>,
  );

  Future<Uri> answerDeepLink(String id, List<LinkTarget> targets) async {
    final json =
        (await _api.post(_deepLinkPath(id), {
              'items': [
                for (final t in targets) {'kind': t.kind, 'id': t.id},
              ],
            }))!
            as Map<String, Object?>;
    return Uri.parse(json['formUrl']! as String);
  }

  String _deepLinkPath(String id) =>
      '/v1/lti/deep-links/${Uri.encodeComponent(id)}';

  String _path(String id) => '/v1/lti/contexts/${Uri.encodeComponent(id)}';

  LtiContext _context(Object? json) =>
      LtiContext.fromJson(json! as Map<String, Object?>);
}
