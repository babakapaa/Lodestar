import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'attendance_labels_context.dart';
import 'attendance_mark.dart';
import 'attendance_status.dart';

class MarkMenu extends StatelessWidget {
  const MarkMenu({required this.mark, required this.onChanged, super.key});

  final AttendanceMark mark;
  final ValueChanged<AttendanceStatus?> onChanged;

  @override
  Widget build(BuildContext context) {
    final labels = context.attendanceLabels;
    return PopupMenuButton<AttendanceStatus?>(
      key: ValueKey('mark-${mark.userId}'),
      onSelected: onChanged,
      itemBuilder: (context) => [
        for (final status in AttendanceStatus.values)
          PopupMenuItem(value: status, child: Text(labels.status(status))),
        PopupMenuItem(child: Text(context.l10n.attendanceClear)),
      ],
      child: Padding(
        padding: const EdgeInsets.all(8),
        child: Text(labels.status(mark.status)),
      ),
    );
  }
}
