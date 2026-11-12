import 'package:lodestar/core/auth/web_authenticator.dart';

class FakeAuthenticator implements WebAuthenticator {
  FakeAuthenticator({this.code = 'one-time-code'});

  final String? code;
  final List<Uri> opened = [];

  @override
  Future<Uri?> authenticate(Uri url, String callbackScheme) async {
    opened.add(url);
    final code = this.code;
    if (code == null) {
      return null;
    }
    return Uri.parse(url.queryParameters['redirect_uri']!)
        .replace(queryParameters: {'code': code});
  }
}
