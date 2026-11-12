import 'dart:convert';
import 'dart:math';

import 'package:crypto/crypto.dart';

class Pkce {
  const Pkce._(this.verifier, this.challenge);

  factory Pkce.generate([Random? random]) {
    final source = random ?? Random.secure();
    final bytes = List<int>.generate(32, (_) => source.nextInt(256));
    final verifier = _encode(bytes);
    return Pkce._(verifier, challengeOf(verifier));
  }

  final String verifier;
  final String challenge;

  static String challengeOf(String verifier) =>
      _encode(sha256.convert(ascii.encode(verifier)).bytes);

  static String _encode(List<int> bytes) =>
      base64Url.encode(bytes).replaceAll('=', '');
}
