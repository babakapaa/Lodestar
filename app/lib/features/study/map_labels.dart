import '../../l10n/generated/app_localizations.dart';
import 'map_gap.dart';
import 'map_rules.dart';
import 'search_map.dart';

class MapLabels {
  const MapLabels(this._l10n);

  final AppLocalizations _l10n;

  String rules(MapRules? rules) => rules == null
      ? _l10n.mapRulesNone
      : _l10n.mapRules(rules.minSources, rules.minScholarly);

  String gap(MapGap gap, SearchMap map, MapRules rules) => switch (gap) {
    MapGap.sources => _l10n.mapGapSources(
      map.used.length,
      rules.minSources < 1 ? 1 : rules.minSources,
    ),
    MapGap.scholarly => _l10n.mapGapScholarly(
      map.scholarlyShare,
      rules.minScholarly,
    ),
    MapGap.why => _l10n.mapGapWhy,
    MapGap.path => _l10n.mapGapPath,
  };
}
