import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/files/files_providers.dart';
import '../../core/files/uploaded_file.dart';
import 'chat_file.dart';

class ChatFiles {
  const ChatFiles(this._ref);

  final WidgetRef _ref;

  Future<List<ChatFile>> pick(int limit) async {
    final opened = await _ref.read(fileOpenerProvider).open(const []);
    final transfer = _ref.read(fileTransferProvider);
    return [
      for (final file in opened.take(limit))
        _chatFile(await transfer.upload(file)),
    ];
  }

  static ChatFile _chatFile(UploadedFile file) => ChatFile(
    id: file.id,
    name: file.name,
    contentType: file.contentType,
    size: file.size,
  );

  Future<void> save(ChatFile file) async {
    final bytes = await _ref.read(fileTransferProvider).download(file.id);
    final dot = file.name.lastIndexOf('.');
    final named = dot > 0;
    await _ref
        .read(exportSaverProvider)
        .save(
          named ? file.name.substring(0, dot) : file.name,
          bytes,
          extension: named ? file.name.substring(dot + 1) : 'bin',
        );
  }
}
