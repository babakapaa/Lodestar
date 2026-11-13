import '../../core/auth/auth_state.dart';
import '../../core/navigation/app_routes.dart';
import 'app_tab.dart';

abstract final class AuthRedirect {
  static const _entries = {AppRoutes.connect, AppRoutes.login};

  static String? target(AuthState auth, String location) {
    if (location == AppRoutes.ltiLaunch) {
      return null;
    }
    final gate = switch (auth) {
      Disconnected() => AppRoutes.connect,
      SignedOut() => AppRoutes.login,
      SignedIn() => null,
    };
    if (gate == null) {
      return _entries.contains(location) ? AppTab.today.path : null;
    }
    if (location == gate || location == AppRoutes.settings) {
      return null;
    }
    return gate;
  }
}
