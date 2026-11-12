class ApiException implements Exception {
  const ApiException(this.statusCode, this.code);

  static const unreachable = 'unreachable';
  static const unauthorized = 'unauthorized';
  static const loginFailed = 'login_failed';

  final int statusCode;
  final String code;

  @override
  String toString() => 'ApiException($statusCode, $code)';
}
