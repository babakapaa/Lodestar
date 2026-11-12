import 'package:flutter/foundation.dart';

class AppRedirect {
  const AppRedirect({required this.uri, required this.callbackScheme});

  factory AppRedirect.current() {
    if (kIsWeb) {
      final base = Uri.base;
      return AppRedirect(
        uri: Uri(
          scheme: base.scheme,
          host: base.host,
          port: base.hasPort ? base.port : null,
          path: '/auth.html',
        ),
        callbackScheme: base.scheme,
      );
    }
    return switch (defaultTargetPlatform) {
      TargetPlatform.windows || TargetPlatform.linux => AppRedirect(
        uri: Uri.parse('http://localhost:$desktopPort/auth'),
        callbackScheme: 'http://localhost:$desktopPort',
      ),
      _ => AppRedirect(
        uri: Uri.parse('$appScheme:/auth'),
        callbackScheme: appScheme,
      ),
    };
  }

  static const appScheme = 'io.github.babakapaa.lodestar';
  static const desktopPort = 43823;

  final Uri uri;
  final String callbackScheme;
}
