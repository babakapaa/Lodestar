import 'package:material_ui/material_ui.dart';

import '../../design/tokens/theme_family.dart';
import 'app_language.dart';
import 'appearance_mode.dart';

@immutable
class AppSettings {
  const AppSettings({
    this.family = ThemeFamily.graphite,
    this.appearance = AppearanceMode.system,
    this.language = AppLanguage.system,
  });

  final ThemeFamily family;
  final AppearanceMode appearance;
  final AppLanguage language;

  AppSettings copyWith({
    ThemeFamily? family,
    AppearanceMode? appearance,
    AppLanguage? language,
  }) {
    return AppSettings(
      family: family ?? this.family,
      appearance: appearance ?? this.appearance,
      language: language ?? this.language,
    );
  }

  @override
  bool operator ==(Object other) =>
      other is AppSettings &&
      other.family == family &&
      other.appearance == appearance &&
      other.language == language;

  @override
  int get hashCode => Object.hash(family, appearance, language);
}
