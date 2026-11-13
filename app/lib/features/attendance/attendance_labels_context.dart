import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'attendance_labels.dart';

extension AttendanceLabelsContext on BuildContext {
  AttendanceLabels get attendanceLabels => AttendanceLabels(l10n);
}
