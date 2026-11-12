import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../../core/auth/identity_provider.dart';

final identityProvidersProvider =
    FutureProvider.autoDispose<List<IdentityProvider>>((ref) {
      final address = ref.watch(authProvider.select((state) => state.address));
      if (address == null) {
        return const [];
      }
      return ref.watch(authApiProvider).providers(address);
    });
