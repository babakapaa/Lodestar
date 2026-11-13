import 'package:flutter/foundation.dart';

import '../../core/auth/app_redirect.dart';
import '../../core/navigation/app_routes.dart';

abstract final class AttendLink {
  static String build(String sessionId, String token) {
    final query = {'s': sessionId, 't': token};
    if (kIsWeb) {
      final base = Uri.base;
      return Uri(
        scheme: base.scheme,
        host: base.host,
        port: base.hasPort ? base.port : null,
        path: AppRoutes.attend,
        queryParameters: query,
      ).toString();
    }
    return Uri(
      scheme: AppRedirect.appScheme,
      path: AppRoutes.attend,
      queryParameters: query,
    ).toString();
  }
}
