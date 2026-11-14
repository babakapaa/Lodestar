import '../../design/tokens/accent_choice.dart';
import '../../design/tokens/theme_family.dart';
import 'app_language.dart';
import 'app_settings.dart';
import 'appearance_mode.dart';
import 'preference_value.dart';
import 'preferences_repository.dart';

class PreferenceSync {
  const PreferenceSync(this._repository);

  static const themeKey = 'core.theme';
  static const accentKey = 'core.accent';
  static const appearanceKey = 'core.appearance';
  static const languageKey = 'core.language';

  final PreferencesRepository _repository;

  Future<AppSettings> pull(AppSettings local) async {
    final values = await _repository.values();
    return local.copyWith(
      family: _chosen(ThemeFamily.values, values[themeKey]),
      accent: _chosen(AccentChoice.values, values[accentKey]),
      appearance: _chosen(AppearanceMode.values, values[appearanceKey]),
      language: _chosen(AppLanguage.values, values[languageKey]),
    );
  }

  Future<void> push(AppSettings previous, AppSettings next) async {
    await Future.wait([
      if (previous.family != next.family)
        _repository.set(themeKey, next.family.name),
      if (previous.accent != next.accent)
        _repository.set(accentKey, next.accent.name),
      if (previous.appearance != next.appearance)
        _repository.set(appearanceKey, next.appearance.name),
      if (previous.language != next.language)
        _repository.set(languageKey, next.language.name),
    ]);
  }

  static T? _chosen<T extends Enum>(List<T> values, PreferenceValue? stored) {
    if (stored == null || !stored.set) {
      return null;
    }
    for (final value in values) {
      if (value.name == stored.value) {
        return value;
      }
    }
    return null;
  }
}
