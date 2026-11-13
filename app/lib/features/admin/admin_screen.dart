import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/access/access_providers.dart';
import '../../core/access/permission_set.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';

class AdminScreen extends ConsumerWidget {
  const AdminScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final permissions =
        ref.watch(myPermissionsProvider).value ?? PermissionSet.empty;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.adminTitle)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: ListView(
              padding: const EdgeInsets.all(Spacing.l),
              children: [
                if (permissions.has(PermissionSet.manageSettings))
                  ListTile(
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.apartment_outlined),
                    title: Text(l10n.adminSettingsTitle),
                    subtitle: Text(l10n.adminSettingsBody),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push(AppRoutes.adminSettings),
                  ),
                if (permissions.has(PermissionSet.manageRoles)) ...[
                  ListTile(
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.badge_outlined),
                    title: Text(l10n.rolesTitle),
                    subtitle: Text(l10n.adminRolesBody),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push(AppRoutes.roles),
                  ),
                  ListTile(
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.people_outline),
                    title: Text(l10n.peopleRolesTitle),
                    subtitle: Text(l10n.adminPeopleBody),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push(AppRoutes.people),
                  ),
                ],
                if (permissions.has(PermissionSet.manageStructure) ||
                    permissions.has(PermissionSet.manageMembers))
                  ListTile(
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.account_tree_outlined),
                    title: Text(l10n.unitsTitle),
                    subtitle: Text(l10n.adminUnitsBody),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push(AppRoutes.units),
                  ),
                if (permissions.has(PermissionSet.viewAudit))
                  ListTile(
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.history),
                    title: Text(l10n.auditTitle),
                    subtitle: Text(l10n.auditBody),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push(AppRoutes.adminAudit),
                  ),
                if (permissions.has(PermissionSet.reviewChatReports))
                  ListTile(
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.flag_outlined),
                    title: Text(l10n.chatReportsTitle),
                    subtitle: Text(l10n.chatReportsBody),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push(AppRoutes.chatReports),
                  ),
                if (permissions.has(PermissionSet.manageAwards))
                  ListTile(
                    key: const ValueKey('admin-awards'),
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.emoji_events_outlined),
                    title: Text(l10n.awardsTitle),
                    subtitle: Text(l10n.awardsHint),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push(AppRoutes.awards),
                  ),
                if (permissions.has(PermissionSet.manageIntegrations))
                  ListTile(
                    key: const ValueKey('admin-lms'),
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.hub_outlined),
                    title: Text(l10n.lmsTitle),
                    subtitle: Text(l10n.lmsHint),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push(AppRoutes.ltiPlatforms),
                  ),
                if (permissions.has(PermissionSet.manageSandbox))
                  ListTile(
                    key: const ValueKey('admin-sandbox'),
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.memory),
                    title: Text(l10n.sandboxAdminTitle),
                    subtitle: Text(l10n.sandboxAdminTile),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push(AppRoutes.sandboxAdmin),
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
