import '../../l10n/generated/app_localizations.dart';

class AuditLabels {
  const AuditLabels(this.l10n);

  final AppLocalizations l10n;

  String action(String action) => switch (action) {
    'core.role.created' => l10n.auditRoleCreated,
    'core.role.updated' => l10n.auditRoleUpdated,
    'core.role.deleted' => l10n.auditRoleDeleted,
    'core.role.moved' => l10n.auditRoleMoved,
    'core.role.assigned' => l10n.auditRoleAssigned,
    'core.role.unassigned' => l10n.auditRoleUnassigned,
    'core.override.set' => l10n.auditOverrideSet,
    'core.setting.changed' => l10n.auditSettingChanged,
    'core.setting.reset' => l10n.auditSettingReset,
    'core.module.enabled' => l10n.auditModuleEnabled,
    'core.module.disabled' => l10n.auditModuleDisabled,
    'core.account.deleted' => l10n.auditAccountDeleted,
    _ => action,
  };
}
