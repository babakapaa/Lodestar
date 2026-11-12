class CrdtException implements Exception {
  const CrdtException(this.code);

  static const badInput = 1;
  static const badRange = 2;

  final int code;

  @override
  String toString() => 'CrdtException($code)';
}
