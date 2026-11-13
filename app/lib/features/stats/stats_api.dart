import '../../core/api/api_client.dart';
import 'stats_entry.dart';
import 'stats_query.dart';
import 'stats_report.dart';

class StatsApi {
  const StatsApi(this._api);

  final ApiClient _api;

  Future<StatsReport> units(StatsQuery query) async => StatsReport.fromJson(
    (await _api.get('/v1/stats/units', {
          if (query.parentId != null) 'parent': query.parentId!,
          'days': '${query.days}',
        }))!
        as Map<String, Object?>,
  );

  Future<List<StatsEntry>> entries() async => [
    for (final item in (await _api.get('/v1/stats/entries'))! as List<Object?>)
      StatsEntry.fromJson(item! as Map<String, Object?>),
  ];
}
