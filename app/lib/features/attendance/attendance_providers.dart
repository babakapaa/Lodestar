import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../admin/admin_providers.dart';
import 'attendance_api.dart';
import 'attendance_session.dart';
import 'course_journal.dart';
import 'open_lesson.dart';

final attendanceApiProvider = Provider<AttendanceApi>(
  (ref) => AttendanceApi(ref.watch(apiClientProvider)),
);

final attendanceEnabledProvider = Provider.autoDispose<bool>(
  (ref) =>
      ref
          .watch(modulesProvider)
          .value
          ?.any((m) => m.name == 'attendance' && m.enabled) ??
      false,
);

final openLessonsProvider = FutureProvider.autoDispose<List<OpenLesson>>(
  (ref) => ref.watch(attendanceApiProvider).openLessons(),
);

final attendanceSessionProvider = FutureProvider.autoDispose
    .family<AttendanceSession, String>(
      (ref, id) => ref.watch(attendanceApiProvider).session(id),
    );

final courseJournalProvider = FutureProvider.autoDispose
    .family<CourseJournal, String>(
      (ref, id) => ref.watch(attendanceApiProvider).journal(id),
    );
