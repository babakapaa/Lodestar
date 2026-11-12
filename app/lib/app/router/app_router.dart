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
import '../../features/notes/note_members_screen.dart';
import '../../features/notes/note_screen.dart';
import '../../features/notes/notes_screen.dart';
import '../../features/notes/notes_trash_screen.dart';
import '../../features/notifications/notification_settings_screen.dart';
import '../../features/notifications/notifications_screen.dart';
import '../../features/roles/people_roles_screen.dart';
import '../../features/roles/person_roles_screen.dart';
import '../../features/privacy/privacy_screen.dart';
import '../../features/profile/profile_edit_screen.dart';
import '../../features/profile/profile_screen.dart';
import '../../features/roles/role_edit_screen.dart';
import '../../features/roles/roles_screen.dart';
import '../../features/settings/settings_screen.dart';
import '../../features/study/assignment_edit_screen.dart';
import '../../features/study/assignment_screen.dart';
import '../../features/study/course_screen.dart';
import '../../features/study/grade_screen.dart';
import '../../features/study/schedule_screen.dart';
import '../../features/study/study_screen.dart';
import '../../features/today/today_screen.dart';
import '../../features/units/units_screen.dart';
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
          _route(AppRoutes.roles, const RolesScreen()),
          _route(AppRoutes.roleNew, const RoleEditScreen()),
          _param(AppRoutes.role(':id'), (id) => RoleEditScreen(roleId: id)),
          _route(AppRoutes.people, const PeopleRolesScreen()),
          _param(
            AppRoutes.personRoles(':id'),
            (id) => PersonRolesScreen(userId: id),
          ),
          _route(AppRoutes.units, const UnitsScreen()),
          _param(AppRoutes.unit(':id'), (id) => UnitsScreen(unitId: id)),
          _param(AppRoutes.course(':id'), (id) => CourseScreen(courseId: id)),
          _param(
            AppRoutes.assignmentNew(':id'),
            (id) => AssignmentEditScreen(courseId: id),
          ),
          _param(
            AppRoutes.assignment(':id'),
            (id) => AssignmentScreen(assignmentId: id),
          ),
          _param(
            AppRoutes.assignmentEdit(':id'),
            (id) => AssignmentEditScreen(assignmentId: id),
          ),
          GoRoute(
            path: AppRoutes.grade(':id', ':student'),
            pageBuilder: (context, state) => _fade(
              state,
              GradeScreen(
                assignmentId: state.pathParameters['id']!,
                studentId: state.pathParameters['student']!,
              ),
            ),
          ),
          _route(AppRoutes.notesTrash, const NotesTrashScreen()),
          _param(AppRoutes.note(':id'), (id) => NoteScreen(noteId: id)),
          _param(
            AppRoutes.noteMembers(':id'),
            (id) => NoteMembersScreen(noteId: id),
          ),
          GoRoute(
            path: AppRoutes.schedulePath,
            pageBuilder: (context, state) => _fade(
              state,
              ScheduleScreen(date: state.uri.queryParameters['date']),
            ),
          ),
        ],
      );

  final GoRouter config;

  static GoRoute _route(String path, Widget screen) {
    return GoRoute(
      path: path,
      pageBuilder: (context, state) => _fade(state, screen),
    );
  }

  static GoRoute _param(String path, Widget Function(String id) screen) {
    return GoRoute(
      path: path,
      pageBuilder: (context, state) =>
          _fade(state, screen(state.pathParameters['id']!)),
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
