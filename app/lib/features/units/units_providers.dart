import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'unit.dart';
import 'unit_details.dart';
import 'unit_member.dart';
import 'units_query.dart';
import 'units_repository.dart';

final unitsRepositoryProvider = Provider<UnitsRepository>(
  (ref) => UnitsRepository(ref.watch(apiClientProvider)),
);

final unitChildrenProvider = FutureProvider.autoDispose
    .family<List<Unit>, UnitsQuery>(
      (ref, query) => ref
          .watch(unitsRepositoryProvider)
          .children(query.parentId, archived: query.archived),
    );

final unitDetailsProvider = FutureProvider.autoDispose
    .family<UnitDetails, String>(
      (ref, id) => ref.watch(unitsRepositoryProvider).details(id),
    );

final unitMembersProvider = FutureProvider.autoDispose
    .family<List<UnitMember>, String>(
      (ref, id) => ref.watch(unitsRepositoryProvider).members(id),
    );
