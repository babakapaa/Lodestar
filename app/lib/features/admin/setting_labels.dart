import '../../core/settings/app_language.dart';
import '../../design/tokens/theme_family.dart';
import '../../l10n/generated/app_localizations.dart';
import '../settings/settings_labels.dart';

class SettingLabels {
  const SettingLabels(this.l10n);

  final AppLocalizations l10n;

  String name(String key) => switch (key) {
    'core.default_language' => l10n.settingDefaultLanguage,
    'core.default_theme' => l10n.settingDefaultTheme,
    'core.time_zone' => l10n.settingTimeZone,
    'core.files.max_file_mb' => l10n.settingMaxFile,
    'core.audit_retention_days' => l10n.settingAuditRetention,
    'core.self_deletion' => l10n.settingSelfDeletion,
    'core.deletion_grace_days' => l10n.settingDeletionGrace,
    _ => key,
  };

  String option(String key, String option) {
    final labels = SettingsLabels(l10n);
    for (final language in AppLanguage.values) {
      if (key == 'core.default_language' && language.name == option) {
        return labels.language(language);
      }
    }
    for (final family in ThemeFamily.values) {
      if (key == 'core.default_theme' && family.name == option) {
        return labels.family(family);
      }
    }
    return option;
  }
}
