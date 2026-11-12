import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/auth/pkce.dart';

void main() {
  test('challenge is unpadded base64url of sha256', () {
    expect(
      Pkce.challengeOf('dBjftJeZ4CVP-mJ92K9OGldzW3SH7e_O0HlxSqPN3mY'),
      '74u57d10p5t3rhunr2ceU80kuw2lwZ2Wvmor9qSgXvI',
    );
  });

  test('every login gets a fresh verifier', () {
    final first = Pkce.generate();
    final second = Pkce.generate();

    expect(first.verifier, isNot(second.verifier));
    expect(first.verifier.length, 43);
    expect(first.challenge, Pkce.challengeOf(first.verifier));
  });
}
