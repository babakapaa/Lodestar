import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'league_providers.dart';

extension LeagueRefresh on WidgetRef {
  void refreshLeague(String kind, String id) =>
      invalidate(leagueProvider((kind, id)));
}
