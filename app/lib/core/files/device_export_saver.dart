import 'dart:typed_data';

import 'package:file_saver/file_saver.dart';

import 'export_saver.dart';

class DeviceExportSaver implements ExportSaver {
  const DeviceExportSaver();

  @override
  Future<void> save(String name, Uint8List bytes) async {
    await FileSaver.instance.saveFile(
      name: name,
      bytes: bytes,
      fileExtension: 'zip',
      mimeType: MimeType.zip,
    );
  }
}
