import 'dart:typed_data';

import 'api_exception.dart';
import 'api_transport.dart';
import 'auth_api.dart';
import 'server_address.dart';
import 'token_pair.dart';
import 'token_store.dart';

class ApiClient {
  ApiClient({
    required this._transport,
    required this._auth,
    required this._address,
    required this._tokens,
    required this._onExpired,
  });

  final ApiTransport _transport;
  final AuthApi _auth;
  final ServerAddress _address;
  final TokenStore _tokens;
  final Future<void> Function() _onExpired;
  Future<TokenPair?>? _refreshing;

  Future<Object?> get(String path, [Map<String, String>? query]) =>
      _send('GET', path, query: query);

  Future<Object?> post(String path, [Object? body]) =>
      _send('POST', path, body: body);

  Future<Object?> put(String path, Object? body) =>
      _send('PUT', path, body: body);

  Future<Object?> delete(String path) => _send('DELETE', path);

  Future<Uint8List> download(String path) => _authorized(
    (token) => _transport.bytes(_address.resolve(path), accessToken: token),
  );

  Future<Object?> _send(
    String method,
    String path, {
    Object? body,
    Map<String, String>? query,
  }) {
    final uri = _address.resolve(path, query);
    return _authorized(
      (token) => _transport.send(method, uri, body: body, accessToken: token),
    );
  }

  Future<T> _authorized<T>(Future<T> Function(String token) call) async {
    final tokens = _tokens.current;
    if (tokens == null) {
      throw const ApiException(401, ApiException.unauthorized);
    }
    try {
      return await call(tokens.accessToken);
    } on ApiException catch (error) {
      if (error.statusCode != 401) {
        rethrow;
      }
    }
    final fresh = await _refresh(tokens);
    if (fresh == null) {
      await _onExpired();
      throw const ApiException(401, ApiException.unauthorized);
    }
    return call(fresh.accessToken);
  }

  Future<TokenPair?> _refresh(TokenPair stale) {
    return _refreshing ??= _renew(stale).whenComplete(() => _refreshing = null);
  }

  Future<TokenPair?> _renew(TokenPair stale) async {
    final current = _tokens.current;
    if (current == null) {
      return null;
    }
    if (current.accessToken != stale.accessToken) {
      return current;
    }
    try {
      final fresh = await _auth.refresh(_address, current);
      await _tokens.save(fresh);
      return fresh;
    } on ApiException catch (error) {
      if (error.statusCode == 400 || error.statusCode == 401) {
        return null;
      }
      rethrow;
    }
  }
}
