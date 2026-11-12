import 'package:file_picker/file_picker.dart';

import 'file_opener.dart';
import 'opened_file.dart';

class DeviceFileOpener implements FileOpener {
  const DeviceFileOpener();

  @override
  Future<List<OpenedFile>> open(List<String> extensions) async {
    final files = await FilePicker.pickFiles(
      type: FileType.custom,
      allowedExtensions: extensions,
    );
    return [
      for (final file in files) OpenedFile(file.name, await file.readAsBytes()),
    ];
  }
}
