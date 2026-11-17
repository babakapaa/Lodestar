import 'package:material_ui/material_ui.dart';

import '../../design/tokens/accent_choice.dart';
import '../../design/tokens/theme_family.dart';
import 'app_language.dart';
import 'appearance_mode.dart';

@immutable
class AppSettings {
  const AppSettings({
    this.family = ThemeFamily.campus,
    this.accent = AccentChoice.auto,
    this.appearance = AppearanceMode.system,
    this.language = AppLanguage.system,
  });

  final ThemeFamily family;
  final AccentChoice accent;
  final AppearanceMode appearance;
  final AppLanguage language;

  AppSettings copyWith({
    ThemeFamily? family,
    AccentChoice? accent,
    AppearanceMode? appearance,
    AppLanguage? language,
  }) {
    return AppSettings(
      family: family ?? this.family,
      accent: accent ?? this.accent,
      appearance: appearance ?? this.appearance,
      language: language ?? this.language,
    );
  }

  @override
  bool operator ==(Object other) =>
      other is AppSettings &&
      other.family == family &&
      other.accent == accent &&
      other.appearance == appearance &&
      other.language == language;

  @override
  int get hashCode => Object.hash(family, accent, appearance, language);
}
