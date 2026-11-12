import 'package:flutter/foundation.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/auth/auth_state.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/tokens/motion.dart';
import '../../features/admin/admin_screen.dart';
import '../../features/admin/audit_screen.dart';
import '../../features/admin/tenant_settings_screen.dart';
import '../../features/chats/chats_screen.dart';
import '../../features/connect/connect_screen.dart';
import '../../features/login/login_screen.dart';
import '../../features/navigator/navigator_screen.dart';
import '../../features/notes/notes_screen.dart';
import '../../features/notifications/notification_settings_screen.dart';
import '../../features/notifications/notifications_screen.dart';
import '../../features/privacy/privacy_screen.dart';
import '../../features/profile/profile_edit_screen.dart';
import '../../features/profile/profile_screen.dart';
import '../../features/settings/settings_screen.dart';
import '../../features/study/study_screen.dart';
import '../../features/today/today_screen.dart';
import '../shell/adaptive_shell.dart';
import 'app_tab.dart';
import 'auth_redirect.dart';

class AppRouter {
  AppRouter(ValueListenable<AuthState> auth)
    : config = GoRouter(
        initialLocation: AppTab.today.path,
        refreshListenable: auth,
        redirect: (context, state) =>
            AuthRedirect.target(auth.value, state.matchedLocation),
        routes: [
          StatefulShellRoute.indexedStack(
            builder: (context, state, shell) => AdaptiveShell(shell: shell),
            branches: [
              _branch(AppTab.today, const TodayScreen()),
              _branch(AppTab.study, const StudyScreen()),
              _branch(AppTab.navigator, const NavigatorScreen()),
              _branch(AppTab.notes, const NotesScreen()),
              _branch(AppTab.chats, const ChatsScreen()),
            ],
          ),
          _route(AppRoutes.connect, const ConnectScreen()),
          _route(AppRoutes.login, const LoginScreen()),
          _route(AppRoutes.profile, const ProfileScreen()),
          _route(AppRoutes.profileEdit, const ProfileEditScreen()),
          _route(AppRoutes.privacy, const PrivacyScreen()),
          _route(AppRoutes.notifications, const NotificationsScreen()),
          _route(
            AppRoutes.notificationSettings,
            const NotificationSettingsScreen(),
          ),
          _route(AppRoutes.settings, const SettingsScreen()),
          _route(AppRoutes.admin, const AdminScreen()),
          _route(AppRoutes.adminSettings, const TenantSettingsScreen()),
          _route(AppRoutes.adminAudit, const AuditScreen()),
        ],
      );

  final GoRouter config;

  static GoRoute _route(String path, Widget screen) {
    return GoRoute(
      path: path,
      pageBuilder: (context, state) => _fade(state, screen),
    );
  }

  static StatefulShellBranch _branch(AppTab tab, Widget screen) {
    return StatefulShellBranch(routes: [_route(tab.path, screen)]);
  }

  static Page<void> _fade(GoRouterState state, Widget child) {
    return CustomTransitionPage<void>(
      key: state.pageKey,
      child: child,
      transitionDuration: Motion.normal,
      reverseTransitionDuration: Motion.fast,
      transitionsBuilder: (context, animation, secondary, child) =>
          FadeTransition(
            opacity: CurvedAnimation(parent: animation, curve: Motion.curve),
            child: child,
          ),
    );
  }
}
