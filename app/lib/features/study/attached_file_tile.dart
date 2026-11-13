import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/files/file_download.dart';
import '../../l10n/l10n_context.dart';
import 'attached_file.dart';

class AttachedFileTile extends ConsumerWidget {
  const AttachedFileTile({required this.file, super.key});

  final AttachedFile file;

  @override
  Widget build(BuildContext context, WidgetRef ref) => ListTile(
    key: ValueKey('file-${file.id}'),
    contentPadding: EdgeInsets.zero,
    leading: const Icon(Icons.attach_file),
    title: Text(file.name),
    trailing: const Icon(Icons.download_outlined),
    onTap: () async {
      final l10n = context.l10n;
      final messenger = ScaffoldMessenger.of(context);
      try {
        await FileDownload(ref).save(file.id, file.name);
      } on Object catch (error) {
        messenger.showSnackBar(
          SnackBar(content: Text(ApiErrorText.of(l10n, error))),
        );
      }
    },
  );
}
