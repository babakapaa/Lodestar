import 'map_gap.dart';
import 'map_rules.dart';
import 'search_map.dart';

class SearchMapDraft {
  const SearchMapDraft({
    required this.map,
    required this.rules,
    required this.gaps,
    required this.changedAt,
  });

  factory SearchMapDraft.fromJson(Map<String, Object?> json) {
    final changedAt = json['changedAt'];
    return SearchMapDraft(
      map: SearchMap.fromJson(json['map']! as Map<String, Object?>),
      rules: MapRules.maybe(json['rules']),
      gaps: MapGap.list(json['missing']),
      changedAt: changedAt is String ? DateTime.parse(changedAt) : null,
    );
  }

  final SearchMap map;
  final MapRules? rules;
  final List<MapGap> gaps;
  final DateTime? changedAt;
}
