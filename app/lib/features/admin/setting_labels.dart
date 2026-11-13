import '../../core/settings/app_language.dart';
import '../../design/tokens/theme_family.dart';
import '../../l10n/generated/app_localizations.dart';
import '../settings/settings_labels.dart';
import 'setting_group.dart';

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
    'study.teachers_edit_schedule' => l10n.settingTeachersEditSchedule,
    'study.defense_recording' => l10n.settingDefenseRecording,
    'chat.retention_days' => l10n.settingChatRetention,
    'chat.calls' => l10n.settingChatCalls,
    'chat.screen_share' => l10n.settingChatScreenShare,
    'chat.call_recording' => l10n.settingChatRecording,
    'core.lti_assignment_source' => l10n.settingLtiSource,
    'sandbox.base_hours' => l10n.settingSandboxBaseHours,
    'sandbox.deadline_days' => l10n.settingSandboxDeadlineDays,
    'sandbox.watch_processes' => l10n.settingSandboxWatch,
    'sandbox.night_discount' => l10n.settingSandboxNightDiscount,
    'sandbox.night_from' => l10n.settingSandboxNightFrom,
    'sandbox.night_to' => l10n.settingSandboxNightTo,
    'attendance.qr' => l10n.settingAttendanceQr,
    'attendance.code' => l10n.settingAttendanceCode,
    'attendance.networks' => l10n.settingAttendanceNetworks,
    'attendance.late_minutes' => l10n.settingAttendanceLate,
    'attendance.absent_alert' => l10n.settingAttendanceAlert,
    'navigator.support_contacts' => l10n.settingNavigatorSupport,
    'navigator.memory_messages' => l10n.settingNavigatorMemory,
    'league.division_gpu_hours' => l10n.settingLeagueGpu,
    _ => key,
  };

  String group(SettingGroup group) => switch (group) {
    SettingGroup.general => l10n.settingGroupGeneral,
    SettingGroup.privacy => l10n.settingGroupPrivacy,
    SettingGroup.study => l10n.settingGroupStudy,
    SettingGroup.navigator => l10n.settingGroupNavigator,
    SettingGroup.attendance => l10n.settingGroupAttendance,
    SettingGroup.sandbox => l10n.settingGroupSandbox,
    SettingGroup.other => l10n.settingGroupOther,
  };

  String option(String key, String option) {
    if (key == 'chat.call_recording') {
      return switch (option) {
        'off' => l10n.settingRecordingOff,
        'admins' => l10n.settingScreenShareAdmins,
        'everyone' => l10n.settingScreenShareEveryone,
        _ => option,
      };
    }
    if (key == 'chat.screen_share') {
      return switch (option) {
        'everyone' => l10n.settingScreenShareEveryone,
        'admins' => l10n.settingScreenShareAdmins,
        _ => option,
      };
    }
    if (key == 'core.lti_assignment_source') {
      return switch (option) {
        'lodestar' => l10n.settingLtiSourceLodestar,
        'lms' => l10n.settingLtiSourceLms,
        'units' => l10n.settingLtiSourceUnits,
        _ => option,
      };
    }
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
