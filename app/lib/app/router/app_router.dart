import 'package:flutter/foundation.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/auth/auth_state.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/tokens/motion.dart';
import '../../features/admin/admin_screen.dart';
import '../../features/admin/audit_screen.dart';
import '../../features/admin/tenant_settings_screen.dart';
import '../../features/chats/chat_info_screen.dart';
import '../../features/chats/chat_reports_screen.dart';
import '../../features/chats/chat_screen.dart';
import '../../features/chats/chats_screen.dart';
import '../../features/chats/course_chats_screen.dart';
import '../../features/connect/connect_screen.dart';
import '../../features/defense/defense_booking_screen.dart';
import '../../features/defense/defense_screen.dart';
import '../../features/defense/defenses_screen.dart';
import '../../features/league/awards_screen.dart';
import '../../features/league/league_player_screen.dart';
import '../../features/league/league_rules_screen.dart';
import '../../features/league/league_screen.dart';
import '../../features/league/league_task_edit_screen.dart';
import '../../features/league/league_task_screen.dart';
import '../../features/login/login_screen.dart';
import '../../features/lti/course_lms_screen.dart';
import '../../features/lti/lms_platforms_screen.dart';
import '../../features/lti/lti_context_screen.dart';
import '../../features/lti/lti_launch.dart';
import '../../features/lti/lti_launch_screen.dart';
import '../../features/navigator/dialog_scope.dart';
import '../../features/navigator/dialog_screen.dart';
import '../../features/navigator/journal_screen.dart';
import '../../features/navigator/material_form_screen.dart';
import '../../features/navigator/materials_screen.dart';
import '../../features/navigator/navigator_screen.dart';
import '../../features/notes/note_members_screen.dart';
import '../../features/notes/note_screen.dart';
import '../../features/notes/note_version_screen.dart';
import '../../features/notes/note_versions_screen.dart';
import '../../features/notes/notes_graph_screen.dart';
import '../../features/notes/notes_screen.dart';
import '../../features/notes/notes_trash_screen.dart';
import '../../core/sources/source.dart';
import '../../features/notes/source_edit_screen.dart';
import '../../features/notes/source_screen.dart';
import '../../features/notes/sources_screen.dart';
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
import '../../features/study/search_map_screen.dart';
import '../../features/study/study_screen.dart';
import '../../features/today/today_screen.dart';
import '../../features/team/team_join_screen.dart';
import '../../features/team/team_profile_screen.dart';
import '../../features/team/team_scope_screen.dart';
import '../../features/team/team_screen.dart';
import '../../features/team/teams_screen.dart';
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
          _route(AppRoutes.chatReports, const ChatReportsScreen()),
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
          _param(
            AppRoutes.searchMap(':id'),
            (id) => SearchMapScreen(assignmentId: id),
          ),
          _param(
            AppRoutes.defenseBooking(':id'),
            (id) => DefenseBookingScreen(assignmentId: id),
          ),
          _param(
            AppRoutes.defenses(':id'),
            (id) => DefensesScreen(assignmentId: id),
          ),
          GoRoute(
            path: AppRoutes.defense(':id', ':student'),
            pageBuilder: (context, state) => _fade(
              state,
              DefenseScreen(
                assignmentId: state.pathParameters['id']!,
                studentId: state.pathParameters['student']!,
              ),
            ),
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
          _route(AppRoutes.notesGraph, const NotesGraphScreen()),
          _route(AppRoutes.sourcesPath, const SourcesScreen()),
          GoRoute(
            path: AppRoutes.sourceNew,
            pageBuilder: (context, state) =>
                _fade(state, SourceEditScreen(source: state.extra as Source?)),
          ),
          GoRoute(
            path: AppRoutes.source(':id'),
            pageBuilder: (context, state) => _fade(
              state,
              SourceScreen(
                sourceId: state.pathParameters['id']!,
                source: state.extra as Source?,
              ),
            ),
          ),
          _param(AppRoutes.note(':id'), (id) => NoteScreen(noteId: id)),
          _param(
            AppRoutes.noteMembers(':id'),
            (id) => NoteMembersScreen(noteId: id),
          ),
          _param(
            AppRoutes.noteVersions(':id'),
            (id) => NoteVersionsScreen(noteId: id),
          ),
          GoRoute(
            path: AppRoutes.noteVersion(':id', ':version'),
            pageBuilder: (context, state) => _fade(
              state,
              NoteVersionScreen(
                noteId: state.pathParameters['id']!,
                versionId: state.pathParameters['version']!,
              ),
            ),
          ),
          GoRoute(
            path: AppRoutes.dialogNew,
            pageBuilder: (context, state) {
              final query = state.uri.queryParameters;
              final kind = query['scopeKind'];
              final id = query['scopeId'];
              return _fade(
                state,
                DialogScreen(
                  scope: kind == null || id == null
                      ? null
                      : DialogScope(kind: kind, id: id),
                ),
              );
            },
          ),
          _param(AppRoutes.dialog(':id'), (id) => DialogScreen(dialogId: id)),
          _route(AppRoutes.journal, const JournalScreen()),
          _param(AppRoutes.chat(':id'), (id) => ChatScreen(chatId: id)),
          GoRoute(
            path: '${AppRoutes.courseChatsPath}/:id',
            pageBuilder: (context, state) => _fade(
              state,
              CourseChatsScreen(
                courseId: state.pathParameters['id']!,
                teaches: state.uri.queryParameters['teach'] == '1',
              ),
            ),
          ),
          _param(AppRoutes.chatInfo(':id'), (id) => ChatInfoScreen(chatId: id)),
          _param(
            AppRoutes.courseMaterials(':id'),
            (id) => MaterialsScreen(courseId: id),
          ),
          _param(
            AppRoutes.materialNew(':id'),
            (id) => MaterialFormScreen(courseId: id),
          ),
          _route(AppRoutes.awards, const AwardsScreen()),
          _param(
            AppRoutes.leagueTask(':id'),
            (id) => LeagueTaskScreen(taskId: id),
          ),
          _param(
            AppRoutes.leagueTaskEdit(':id'),
            (id) => LeagueTaskEditScreen(taskId: id),
          ),
          _arena(
            AppRoutes.league(':kind', ':id'),
            (kind, id) => LeagueScreen(kind: kind, originId: id),
          ),
          _arena(
            AppRoutes.leagueRules(':kind', ':id'),
            (kind, id) => LeagueRulesScreen(kind: kind, originId: id),
          ),
          _arena(
            AppRoutes.leagueTaskNew(':kind', ':id'),
            (kind, id) => LeagueTaskEditScreen(kind: kind, originId: id),
          ),
          GoRoute(
            path: AppRoutes.leaguePlayer(':kind', ':id', ':user'),
            pageBuilder: (context, state) => _fade(
              state,
              LeaguePlayerScreen(
                kind: state.pathParameters['kind']!,
                originId: state.pathParameters['id']!,
                userId: state.pathParameters['user']!,
              ),
            ),
          ),
          _route(AppRoutes.teams, const TeamsScreen()),
          _route(AppRoutes.teamProfile, const TeamProfileScreen()),
          _param(
            AppRoutes.teamJoin(':id'),
            (code) => TeamJoinScreen(code: code),
          ),
          GoRoute(
            path: AppRoutes.teamScope(':kind', ':id'),
            pageBuilder: (context, state) => _fade(
              state,
              TeamScopeScreen(
                kind: state.pathParameters['kind']!,
                scopeId: state.pathParameters['id']!,
              ),
            ),
          ),
          _param(AppRoutes.team(':id'), (id) => TeamScreen(teamId: id)),
          GoRoute(
            path: AppRoutes.ltiLaunch,
            pageBuilder: (context, state) => _fade(
              state,
              LtiLaunchScreen(
                launch: LtiLaunch.parse(state.uri.queryParameters),
              ),
            ),
          ),
          _param(
            AppRoutes.ltiContext(':id'),
            (id) => LtiContextScreen(contextId: id),
          ),
          _route(AppRoutes.ltiPlatforms, const LmsPlatformsScreen()),
          _param(
            AppRoutes.courseLms(':id'),
            (id) => CourseLmsScreen(courseId: id),
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

  static GoRoute _arena(
    String path,
    Widget Function(String kind, String id) screen,
  ) {
    return GoRoute(
      path: path,
      pageBuilder: (context, state) => _fade(
        state,
        screen(state.pathParameters['kind']!, state.pathParameters['id']!),
      ),
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
