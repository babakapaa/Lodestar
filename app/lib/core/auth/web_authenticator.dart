abstract interface class WebAuthenticator {
  Future<Uri?> authenticate(Uri url, String callbackScheme);
}
