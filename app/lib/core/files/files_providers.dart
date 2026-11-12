import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'device_export_saver.dart';
import 'device_file_opener.dart';
import 'export_saver.dart';
import 'file_opener.dart';

final exportSaverProvider = Provider<ExportSaver>(
  (ref) => const DeviceExportSaver(),
);

final fileOpenerProvider = Provider<FileOpener>(
  (ref) => const DeviceFileOpener(),
);
