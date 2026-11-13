import 'dart:convert';
import 'dart:typed_data';

import 'package:http/http.dart' as http;

import 'api_exception.dart';

class ApiTransport {
  const ApiTransport(this._client);

  static const _timeout = Duration(seconds: 30);

  final http.Client _client;

  Future<Object?> send(
    String method,
    Uri uri, {
    Object? body,
    String? accessToken,
  }) async {
    final response = await _exchange(method, uri, body, accessToken);
    return _decode(response);
  }

  Future<Uint8List> bytes(Uri uri, {required String accessToken}) async {
    final response = await _exchange('GET', uri, null, accessToken);
    return response.bodyBytes;
  }

  Future<Object?> upload(
    Uri uri,
    Uint8List bytes, {
    required String contentType,
    required String accessToken,
    Set<int> accepted = const {},
  }) async {
    final response = await _exchange(
      'POST',
      uri,
      null,
      accessToken,
      raw: bytes,
      contentType: contentType,
      accepted: accepted,
    );
    return _decode(response);
  }

  Future<http.Response> _exchange(
    String method,
    Uri uri,
    Object? body,
    String? accessToken, {
    Uint8List? raw,
    String? contentType,
    Set<int> accepted = const {},
  }) async {
    final request = http.Request(method, uri)
      ..headers['Accept'] = 'application/json';
    if (accessToken != null) {
      request.headers['Authorization'] = 'Bearer $accessToken';
    }
    if (body != null) {
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(body);
    }
    if (raw != null) {
      request.headers['Content-Type'] =
          contentType ?? 'application/octet-stream';
      request.bodyBytes = raw;
    }
    final http.Response response;
    try {
      response = await http.Response.fromStream(
        await _client.send(request).timeout(_timeout),
      );
    } on Exception {
      throw const ApiException(0, ApiException.unreachable);
    }
    if ((response.statusCode >= 200 && response.statusCode < 300) ||
        accepted.contains(response.statusCode)) {
      return response;
    }
    final decoded = _decode(response);
    final code = decoded is Map<String, Object?> ? decoded['code'] : null;
    throw ApiException(response.statusCode, code is String ? code : 'unknown');
  }

  Object? _decode(http.Response response) {
    if (response.bodyBytes.isEmpty) {
      return null;
    }
    try {
      return jsonDecode(utf8.decode(response.bodyBytes));
    } on FormatException {
      return null;
    }
  }
}
