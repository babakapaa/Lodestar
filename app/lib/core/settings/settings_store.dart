import 'package:shared_preferences/shared_preferences.dart';

import '../../design/tokens/accent_choice.dart';
import '../../design/tokens/theme_family.dart';
import 'app_language.dart';
import 'app_settings.dart';
import 'appearance_mode.dart';

class SettingsStore {
  const SettingsStore(this._prefs);

  static const _familyKey = 'settings.family';
  static const _accentKey = 'settings.accent';
  static const _appearanceKey = 'settings.appearance';
  static const _languageKey = 'settings.language';

  final SharedPreferences _prefs;

  AppSettings load() {
    const defaults = AppSettings();
    return AppSettings(
      family: _read(ThemeFamily.values, _familyKey) ?? defaults.family,
      accent: _read(AccentChoice.values, _accentKey) ?? defaults.accent,
      appearance:
          _read(AppearanceMode.values, _appearanceKey) ?? defaults.appearance,
      language: _read(AppLanguage.values, _languageKey) ?? defaults.language,
    );
  }

  Future<void> save(AppSettings settings) async {
    await Future.wait([
      _prefs.setString(_familyKey, settings.family.name),
      _prefs.setString(_accentKey, settings.accent.name),
      _prefs.setString(_appearanceKey, settings.appearance.name),
      _prefs.setString(_languageKey, settings.language.name),
    ]);
  }

  T? _read<T extends Enum>(List<T> values, String key) {
    final name = _prefs.getString(key);
    for (final value in values) {
      if (value.name == name) {
        return value;
      }
    }
    return null;
  }
}
