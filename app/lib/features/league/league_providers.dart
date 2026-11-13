import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'award.dart';
import 'league_api.dart';
import 'league_view.dart';
import 'player_details.dart';
import 'task_details.dart';

final leagueApiProvider = Provider<LeagueApi>(
  (ref) => LeagueApi(ref.watch(apiClientProvider)),
);

final leagueProvider = FutureProvider.autoDispose
    .family<LeagueView, (String, String)>(
      (ref, key) => ref.watch(leagueApiProvider).view(key.$1, key.$2),
    );

final playerProvider = FutureProvider.autoDispose
    .family<PlayerDetails, (String, String)>(
      (ref, key) => ref.watch(leagueApiProvider).player(key.$1, key.$2),
    );

final leagueTaskProvider = FutureProvider.autoDispose
    .family<TaskDetails, String>(
      (ref, id) => ref.watch(leagueApiProvider).task(id),
    );

final awardsProvider = FutureProvider.autoDispose<List<Award>>(
  (ref) => ref.watch(leagueApiProvider).awards(),
);
