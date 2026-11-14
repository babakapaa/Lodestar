import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/access/access_providers.dart';
import '../../core/auth/auth_providers.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/tile_group.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../stats/stats_providers.dart';
import '../tester/tester_section.dart';
import 'profile_header.dart';
import 'profile_providers.dart';
import 'sessions_section.dart';

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final administers =
        ref.watch(myPermissionsProvider).value?.administers ?? false;
    final units = ref.watch(statsEntriesProvider).value ?? const [];
    final email = ref.watch(accountProvider).value?.email;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.profileTitle)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: ListView(
              padding: const EdgeInsets.all(Spacing.l),
              children: [
                const ProfileHeader(),
                const SizedBox(height: Spacing.xl),
                TileGroup(
                  children: [
                    for (final (icon, title, route) in [
                      (Icons.tune, l10n.settingsTitle, AppRoutes.settings),
                      (
                        Icons.notifications_outlined,
                        l10n.notificationSettingsTitle,
                        AppRoutes.notificationSettings,
                      ),
                      (
                        Icons.shield_outlined,
                        l10n.privacyTitle,
                        AppRoutes.privacy,
                      ),
                      if (administers)
                        (
                          Icons.admin_panel_settings_outlined,
                          l10n.adminTitle,
                          AppRoutes.admin,
                        ),
                      for (final unit in units)
                        (
                          Icons.query_stats,
                          l10n.statsOfUnit(unit.name),
                          AppRoutes.unitStats(unit.unitId, unit.name),
                        ),
                    ])
                      ListTile(
                        contentPadding: TileGroup.tilePadding,
                        leading: Icon(icon),
                        title: Text(title),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () => context.push(route),
                      ),
                  ],
                ),
                const SizedBox(height: Spacing.xl),
                const SessionsSection(),
                TesterSection(hint: l10n.testerSwitchHint, current: email),
                const SizedBox(height: Spacing.xl),
                LodestarButton(
                  label: l10n.signOut,
                  icon: Icons.logout,
                  variant: LodestarButtonVariant.secondary,
                  onPressed: () => ref.read(authProvider.notifier).signOut(),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
