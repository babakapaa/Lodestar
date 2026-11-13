import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';

import '../../core/auth/auth_providers.dart';
import 'lesson.dart';
import 'lesson_series.dart';
import 'schedule_api.dart';

const upcomingDays = 14;

final scheduleApiProvider = Provider<ScheduleApi>(
  (ref) => ScheduleApi(ref.watch(apiClientProvider)),
);

final courseSeriesProvider = FutureProvider.autoDispose
    .family<List<LessonSeries>, String>(
      (ref, courseId) => ref.watch(scheduleApiProvider).series(courseId),
    );

final upcomingLessonsProvider = FutureProvider.autoDispose
    .family<List<Lesson>, String>((ref, courseId) {
      final today = DateTime.now();
      final format = DateFormat('yyyy-MM-dd');
      return ref
          .watch(scheduleApiProvider)
          .lessons(
            courseId,
            format.format(today),
            format.format(today.add(const Duration(days: upcomingDays))),
          );
    });
