import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'files_providers.dart';

class FileDownload {
  const FileDownload(this._ref);

  final WidgetRef _ref;

  Future<void> save(String id, String name) async {
    final bytes = await _ref.read(fileTransferProvider).download(id);
    final dot = name.lastIndexOf('.');
    final named = dot > 0;
    await _ref
        .read(exportSaverProvider)
        .save(
          named ? name.substring(0, dot) : name,
          bytes,
          extension: named ? name.substring(dot + 1) : 'bin',
        );
  }
}
