import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../auth/auth_providers.dart';
import '../auth/auth_state.dart';
import 'permission_set.dart';

final myPermissionsProvider = FutureProvider<PermissionSet>((ref) async {
  if (ref.watch(authProvider) is! SignedIn) {
    return PermissionSet.empty;
  }
  final body = await ref.watch(apiClientProvider).get('/v1/me/permissions');
  return PermissionSet({
    for (final item in body! as List<Object?>) item! as String,
  });
});
