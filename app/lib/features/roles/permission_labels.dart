import '../../l10n/generated/app_localizations.dart';

class PermissionLabels {
  const PermissionLabels(this.l10n);

  final AppLocalizations l10n;

  String of(String permission) => switch (permission) {
    'core.administrator' => l10n.permissionAdministrator,
    'core.roles.manage' => l10n.permissionManageRoles,
    'core.settings.manage' => l10n.permissionManageSettings,
    'core.structure.manage' => l10n.permissionManageStructure,
    'core.structure.members' => l10n.permissionManageMembers,
    'core.audit.view' => l10n.permissionViewAudit,
    'core.integrations.manage' => l10n.permissionManageIntegrations,
    _ => permission,
  };
}
