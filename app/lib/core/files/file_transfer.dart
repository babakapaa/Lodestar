import 'dart:typed_data';

import 'package:http/http.dart' as http;

import '../api/api_client.dart';
import '../api/api_exception.dart';
import 'content_types.dart';
import 'opened_file.dart';
import 'uploaded_file.dart';

class FileTransfer {
  const FileTransfer(this._api, this._http);

  final ApiClient _api;
  final http.Client _http;

  Future<UploadedFile> upload(OpenedFile file) async {
    final created =
        (await _api.post('/v1/files', {
              'name': file.name,
              'contentType': ContentTypes.of(file.name),
              'size': file.bytes.length,
            }))!
            as Map<String, Object?>;
    final id = (created['file']! as Map<String, Object?>)['id']! as String;
    final upload = created['upload']! as Map<String, Object?>;
    final request =
        http.Request(
            upload['method'] as String? ?? 'PUT',
            Uri.parse(upload['url']! as String),
          )
          ..headers.addAll({
            for (final MapEntry(:key, :value)
                in (upload['headers'] as Map<String, Object?>? ?? const {})
                    .entries)
              key: value! as String,
          })
          ..bodyBytes = file.bytes;
    final response = await _send(request);
    if (response.statusCode >= 300) {
      throw ApiException(response.statusCode, 'upload_failed');
    }
    return UploadedFile.fromJson(
      (await _api.post('/v1/files/${Uri.encodeComponent(id)}/complete'))!
          as Map<String, Object?>,
    );
  }

  Future<Uint8List> download(String id) async {
    final link =
        (await _api.get('/v1/files/${Uri.encodeComponent(id)}/download'))!
            as Map<String, Object?>;
    final response = await _send(
      http.Request('GET', Uri.parse(link['url']! as String)),
    );
    if (response.statusCode >= 300) {
      throw ApiException(response.statusCode, 'download_failed');
    }
    return response.bodyBytes;
  }

  Future<http.Response> _send(http.Request request) async {
    try {
      return await http.Response.fromStream(await _http.send(request));
    } on http.ClientException {
      throw const ApiException(0, ApiException.unreachable);
    }
  }
}
