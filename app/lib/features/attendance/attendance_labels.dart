import '../../l10n/generated/app_localizations.dart';
import 'attendance_status.dart';

class AttendanceLabels {
  const AttendanceLabels(this.l10n);

  final AppLocalizations l10n;

  String status(AttendanceStatus? status) => switch (status) {
    AttendanceStatus.present => l10n.attendancePresent,
    AttendanceStatus.late => l10n.attendanceLate,
    AttendanceStatus.absent => l10n.attendanceAbsent,
    AttendanceStatus.excused => l10n.attendanceExcused,
    null => l10n.attendanceUnmarked,
  };
}
