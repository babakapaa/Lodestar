import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../auth/auth_providers.dart';

import 'device_export_saver.dart';
import 'device_file_opener.dart';
import 'export_saver.dart';
import 'file_opener.dart';
import 'file_transfer.dart';

final exportSaverProvider = Provider<ExportSaver>(
  (ref) => const DeviceExportSaver(),
);

final fileOpenerProvider = Provider<FileOpener>(
  (ref) => const DeviceFileOpener(),
);

final fileTransferProvider = Provider<FileTransfer>(
  (ref) =>
      FileTransfer(ref.watch(apiClientProvider), ref.watch(httpClientProvider)),
);
