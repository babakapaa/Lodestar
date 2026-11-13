import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/files/files_providers.dart';
import '../../l10n/l10n_context.dart';
import '../study/table_format.dart';
import '../team/team_action.dart';
import 'attendance_providers.dart';

class JournalExportMenu extends ConsumerWidget {
  const JournalExportMenu({required this.courseId, super.key});

  final String courseId;

  Future<void> _export(
    BuildContext context,
    WidgetRef ref,
    TableFormat format,
  ) => TeamAction.run(context, () async {
    final bytes = await ref
        .read(attendanceApiProvider)
        .export(courseId, format);
    final day = DateTime.now().toIso8601String().substring(0, 10);
    await ref
        .read(exportSaverProvider)
        .save('attendance-$day', bytes, extension: format.name);
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return PopupMenuButton<TableFormat>(
      key: const ValueKey('attendance-export'),
      tooltip: context.l10n.attendanceExport,
      icon: const Icon(Icons.download_outlined),
      onSelected: (format) => _export(context, ref, format),
      itemBuilder: (context) => [
        for (final format in TableFormat.values)
          PopupMenuItem(value: format, child: Text(format.name.toUpperCase())),
      ],
    );
  }
}
