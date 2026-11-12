import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../design/tokens/theme_family.dart';
import 'app_language.dart';
import 'app_settings.dart';
import 'appearance_mode.dart';
import 'settings_providers.dart';

class SettingsController extends Notifier<AppSettings> {
  @override
  AppSettings build() => ref.watch(settingsStoreProvider).load();

  Future<void> setFamily(ThemeFamily family) =>
      _update(state.copyWith(family: family));

  Future<void> setAppearance(AppearanceMode appearance) =>
      _update(state.copyWith(appearance: appearance));

  Future<void> setLanguage(AppLanguage language) =>
      _update(state.copyWith(language: language));

  Future<void> _update(AppSettings next) async {
    if (next == state) {
      return;
    }
    state = next;
    await ref.read(settingsStoreProvider).save(next);
  }
}
