import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'team_providers.dart';
import 'team_scope.dart';

extension TeamScopeRefresh on WidgetRef {
  void refreshScope(TeamScope scope) {
    invalidate(teamScopeProvider((scope.kind, scope.id)));
    invalidate(myTeamsProvider);
  }
}
