import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/motion.dart';
import '../../features/chats/chats_screen.dart';
import '../../features/navigator/navigator_screen.dart';
import '../../features/notes/notes_screen.dart';
import '../../features/settings/settings_screen.dart';
import '../../features/study/study_screen.dart';
import '../../features/today/today_screen.dart';
import '../shell/adaptive_shell.dart';
import '../../core/navigation/app_routes.dart';
import 'app_tab.dart';

class AppRouter {
  AppRouter();

  final GoRouter config = GoRouter(
    initialLocation: AppTab.today.path,
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
      GoRoute(
        path: AppRoutes.settings,
        pageBuilder: (context, state) => _fade(state, const SettingsScreen()),
      ),
    ],
  );

  static StatefulShellBranch _branch(AppTab tab, Widget screen) {
    return StatefulShellBranch(
      routes: [
        GoRoute(
          path: tab.path,
          pageBuilder: (context, state) => _fade(state, screen),
        ),
      ],
    );
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
