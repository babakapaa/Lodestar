import '../../core/settings/app_language.dart';
import '../../core/settings/appearance_mode.dart';
import '../../design/tokens/accent_choice.dart';
import '../../design/tokens/theme_family.dart';
import '../../l10n/generated/app_localizations.dart';

class SettingsLabels {
  const SettingsLabels(this.l10n);

  final AppLocalizations l10n;

  String family(ThemeFamily family) => switch (family) {
    ThemeFamily.paper => l10n.themePaper,
    ThemeFamily.graphite => l10n.themeGraphite,
    ThemeFamily.lab => l10n.themeLab,
    ThemeFamily.campus => l10n.themeCampus,
  };

  String accent(AccentChoice accent) => switch (accent) {
    AccentChoice.auto => l10n.accentAuto,
    AccentChoice.blue => l10n.accentBlue,
    AccentChoice.indigo => l10n.accentIndigo,
    AccentChoice.violet => l10n.accentViolet,
    AccentChoice.pink => l10n.accentPink,
    AccentChoice.red => l10n.accentRed,
    AccentChoice.amber => l10n.accentAmber,
    AccentChoice.green => l10n.accentGreen,
    AccentChoice.teal => l10n.accentTeal,
  };

  String appearance(AppearanceMode mode) => switch (mode) {
    AppearanceMode.system => l10n.appearanceSystem,
    AppearanceMode.light => l10n.appearanceLight,
    AppearanceMode.dark => l10n.appearanceDark,
  };

  String language(AppLanguage language) => switch (language) {
    AppLanguage.system => l10n.languageSystem,
    AppLanguage.kk => l10n.languageKazakh,
    AppLanguage.ru => l10n.languageRussian,
    AppLanguage.en => l10n.languageEnglish,
  };
}
