import 'dart:typed_data';

abstract interface class ExportSaver {
  Future<void> save(String name, Uint8List bytes, {String extension = 'zip'});
}
