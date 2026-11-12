import 'token_pair.dart';

abstract interface class TokenStore {
  TokenPair? get current;

  Future<void> save(TokenPair tokens);

  Future<void> clear();
}
