import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';

final actorNameProvider = FutureProvider.autoDispose.family<String, String>((
  ref,
  userId,
) async {
  final body = await ref
      .watch(apiClientProvider)
      .get('/v1/users/${Uri.encodeComponent(userId)}/profile');
  return (body! as Map<String, Object?>)['displayName']! as String;
});
