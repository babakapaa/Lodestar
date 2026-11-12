import 'package:flutter/services.dart';
import 'package:flutter_web_auth_2/flutter_web_auth_2.dart';

import 'web_authenticator.dart';

class BrowserAuthenticator implements WebAuthenticator {
  const BrowserAuthenticator();

  @override
  Future<Uri?> authenticate(Uri url, String callbackScheme) async {
    try {
      final result = await FlutterWebAuth2.authenticate(
        url: url.toString(),
        callbackUrlScheme: callbackScheme,
      );
      return Uri.parse(result);
    } on PlatformException catch (error) {
      if (error.code == 'CANCELED') {
        return null;
      }
      rethrow;
    }
  }
}
