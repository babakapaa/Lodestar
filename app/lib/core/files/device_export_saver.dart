import 'dart:typed_data';

import 'package:file_saver/file_saver.dart';

import 'export_saver.dart';

class DeviceExportSaver implements ExportSaver {
  const DeviceExportSaver();

  @override
  Future<void> save(
    String name,
    Uint8List bytes, {
    String extension = 'zip',
  }) async {
    await FileSaver.instance.saveFile(
      name: name,
      bytes: bytes,
      fileExtension: extension,
      mimeType: switch (extension) {
        'md' => MimeType.markdown,
        'zip' => MimeType.zip,
        _ => MimeType.other,
      },
    );
  }
}
