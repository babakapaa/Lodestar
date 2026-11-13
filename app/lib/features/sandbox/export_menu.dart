import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/files/files_providers.dart';
import '../../l10n/l10n_context.dart';
import '../study/table_format.dart';
import '../team/team_action.dart';
import 'sandbox_providers.dart';

class ExportMenu extends ConsumerWidget {
  const ExportMenu({super.key});

  Future<void> _export(
    BuildContext context,
    WidgetRef ref,
    TableFormat format,
  ) => TeamAction.run(context, () async {
    final bytes = await ref.read(sandboxApiProvider).export(format);
    final day = DateTime.now().toIso8601String().substring(0, 10);
    await ref
        .read(exportSaverProvider)
        .save('gpu-runs-$day', bytes, extension: format.name);
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return PopupMenuButton<TableFormat>(
      key: const ValueKey('sandbox-export'),
      tooltip: l10n.sandboxExport,
      icon: const Icon(Icons.download_outlined),
      onSelected: (format) => _export(context, ref, format),
      itemBuilder: (context) => [
        for (final format in TableFormat.values)
          PopupMenuItem(value: format, child: Text(format.name.toUpperCase())),
      ],
    );
  }
}
