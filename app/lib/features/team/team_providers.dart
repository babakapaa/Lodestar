import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'team_api.dart';
import 'team_card.dart';
import 'team_details.dart';
import 'team_profile.dart';
import 'team_scope_view.dart';

final teamApiProvider = Provider<TeamApi>(
  (ref) => TeamApi(ref.watch(apiClientProvider)),
);

final myTeamsProvider = FutureProvider.autoDispose<List<TeamCard>>(
  (ref) => ref.watch(teamApiProvider).mine(),
);

final teamInvitationsProvider = FutureProvider.autoDispose<List<TeamCard>>(
  (ref) => ref.watch(teamApiProvider).invitations(),
);

final teamScopeProvider = FutureProvider.autoDispose
    .family<TeamScopeView, (String, String)>(
      (ref, key) => ref.watch(teamApiProvider).scope(key.$1, key.$2),
    );

final teamProvider = FutureProvider.autoDispose.family<TeamDetails, String>(
  (ref, id) => ref.watch(teamApiProvider).get(id),
);

final teamPreviewProvider = FutureProvider.autoDispose.family<TeamCard, String>(
  (ref, code) => ref.watch(teamApiProvider).preview(code),
);

final teamProfileProvider = FutureProvider.autoDispose<TeamProfile>(
  (ref) => ref.watch(teamApiProvider).profile(),
);
