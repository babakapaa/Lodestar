import 'dart:typed_data';

class OpenedFile {
  const OpenedFile(this.name, this.bytes);

  final String name;
  final Uint8List bytes;
}
