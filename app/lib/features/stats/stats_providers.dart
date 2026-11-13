import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'stats_api.dart';
import 'stats_query.dart';
import 'stats_report.dart';

final statsApiProvider = Provider<StatsApi>(
  (ref) => StatsApi(ref.watch(apiClientProvider)),
);

final unitStatsProvider = FutureProvider.autoDispose
    .family<StatsReport, StatsQuery>(
      (ref, query) => ref.watch(statsApiProvider).units(query),
    );
