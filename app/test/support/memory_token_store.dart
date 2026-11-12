import 'package:lodestar/core/api/token_pair.dart';
import 'package:lodestar/core/api/token_store.dart';

class MemoryTokenStore implements TokenStore {
  MemoryTokenStore([this._current]);

  TokenPair? _current;

  @override
  TokenPair? get current => _current;

  @override
  Future<void> save(TokenPair tokens) async => _current = tokens;

  @override
  Future<void> clear() async => _current = null;
}
