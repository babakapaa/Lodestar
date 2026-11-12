import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import 'token_pair.dart';
import 'token_store.dart';

class SecureTokenStore implements TokenStore {
  SecureTokenStore._(this._storage, this._current);

  static const _key = 'auth.tokens';

  final FlutterSecureStorage _storage;
  TokenPair? _current;

  static Future<SecureTokenStore> open(FlutterSecureStorage storage) async {
    final raw = await storage.read(key: _key);
    return SecureTokenStore._(storage, _decode(raw));
  }

  @override
  TokenPair? get current => _current;

  @override
  Future<void> save(TokenPair tokens) async {
    _current = tokens;
    await _storage.write(key: _key, value: jsonEncode(tokens.toJson()));
  }

  @override
  Future<void> clear() async {
    _current = null;
    await _storage.delete(key: _key);
  }

  static TokenPair? _decode(String? raw) {
    if (raw == null) {
      return null;
    }
    try {
      return TokenPair.fromJson(jsonDecode(raw) as Map<String, Object?>);
    } on Object {
      return null;
    }
  }
}
