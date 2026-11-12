import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../auth/auth_providers.dart';
import '../auth/auth_state.dart';
import 'app_settings.dart';
import 'preference_sync.dart';
import 'preference_value.dart';
import 'preferences_repository.dart';
import 'settings_controller.dart';
import 'settings_store.dart';

final sharedPreferencesProvider = Provider<SharedPreferences>(
  (ref) => throw StateError(
    'sharedPreferencesProvider must be overridden at startup',
  ),
);

final settingsStoreProvider = Provider<SettingsStore>(
  (ref) => SettingsStore(ref.watch(sharedPreferencesProvider)),
);

final settingsProvider = NotifierProvider<SettingsController, AppSettings>(
  SettingsController.new,
);

final preferencesRepositoryProvider = Provider<PreferencesRepository>(
  (ref) => PreferencesRepository(ref.watch(apiClientProvider)),
);

final preferenceSyncProvider = Provider<PreferenceSync?>((ref) {
  if (ref.watch(authProvider) is! SignedIn) {
    return null;
  }
  return PreferenceSync(ref.watch(preferencesRepositoryProvider));
});

final serverPreferencesProvider =
    FutureProvider.autoDispose<Map<String, PreferenceValue>>(
      (ref) => ref.watch(preferencesRepositoryProvider).values(),
    );
