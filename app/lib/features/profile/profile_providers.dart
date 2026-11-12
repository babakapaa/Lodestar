import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'account.dart';
import 'device_session.dart';
import 'profile.dart';
import 'profile_repository.dart';

final profileRepositoryProvider = Provider<ProfileRepository>(
  (ref) => ProfileRepository(ref.watch(apiClientProvider)),
);

final accountProvider = FutureProvider.autoDispose<Account>(
  (ref) => ref.watch(profileRepositoryProvider).account(),
);

final myProfileProvider = FutureProvider.autoDispose<Profile>(
  (ref) => ref.watch(profileRepositoryProvider).profile(),
);

final sessionsProvider = FutureProvider.autoDispose<List<DeviceSession>>(
  (ref) => ref.watch(profileRepositoryProvider).sessions(),
);
