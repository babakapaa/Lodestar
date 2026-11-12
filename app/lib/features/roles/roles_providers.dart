import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'role.dart';
import 'role_assignment.dart';
import 'roles_repository.dart';

final rolesRepositoryProvider = Provider<RolesRepository>(
  (ref) => RolesRepository(ref.watch(apiClientProvider)),
);

final rolesProvider = FutureProvider.autoDispose<List<Role>>(
  (ref) => ref.watch(rolesRepositoryProvider).roles(),
);

final permissionCatalogProvider = FutureProvider.autoDispose<List<String>>(
  (ref) => ref.watch(rolesRepositoryProvider).catalog(),
);

final personRolesProvider = FutureProvider.autoDispose
    .family<List<RoleAssignment>, String>(
      (ref, userId) => ref.watch(rolesRepositoryProvider).assignmentsOf(userId),
    );
