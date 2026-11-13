import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'team_details.dart';
import 'team_providers.dart';

extension TeamRefresh on WidgetRef {
  void refreshTeam(TeamDetails details) {
    invalidate(teamProvider(details.team.id));
    invalidate(teamScopeProvider((details.scope.kind, details.scope.id)));
    invalidate(myTeamsProvider);
  }
}
