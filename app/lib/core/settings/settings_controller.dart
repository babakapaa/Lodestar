import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../design/tokens/accent_choice.dart';
import '../../design/tokens/theme_family.dart';
import '../api/api_exception.dart';
import 'app_language.dart';
import 'app_settings.dart';
import 'appearance_mode.dart';
import 'settings_providers.dart';

class SettingsController extends Notifier<AppSettings> {
  @override
  AppSettings build() => ref.watch(settingsStoreProvider).load();

  Future<void> setFamily(ThemeFamily family) =>
      _update(state.copyWith(family: family));

  Future<void> setAccent(AccentChoice accent) =>
      _update(state.copyWith(accent: accent));

  Future<void> setAppearance(AppearanceMode appearance) =>
      _update(state.copyWith(appearance: appearance));

  Future<void> setLanguage(AppLanguage language) =>
      _update(state.copyWith(language: language));

  Future<void> pull() async {
    final sync = ref.read(preferenceSyncProvider);
    if (sync == null) {
      return;
    }
    try {
      await _save(await sync.pull(state));
    } on ApiException {
      return;
    }
  }

  Future<void> _update(AppSettings next) async {
    final previous = state;
    if (!await _save(next)) {
      return;
    }
    final sync = ref.read(preferenceSyncProvider);
    if (sync == null) {
      return;
    }
    try {
      await sync.push(previous, next);
    } on ApiException {
      return;
    }
  }

  Future<bool> _save(AppSettings next) async {
    if (next == state) {
      return false;
    }
    state = next;
    await ref.read(settingsStoreProvider).save(next);
    return true;
  }
}
