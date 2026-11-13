import 'dart:typed_data';

import 'package:lodestar/core/files/export_saver.dart';

class FakeExportSaver implements ExportSaver {
  final Map<String, Uint8List> saved = {};

  @override
  Future<void> save(
    String name,
    Uint8List bytes, {
    String extension = 'zip',
  }) async => saved['$name.$extension'] = bytes;
}
