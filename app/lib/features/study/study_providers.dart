import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../../core/offline/loaded.dart';
import '../../core/offline/offline_providers.dart';
import '../../core/settings/settings_providers.dart';
import 'agenda.dart';
import 'assignment_details.dart';
import 'assignment_item.dart';
import 'course.dart';
import 'course_details.dart';
import 'lesson.dart';
import 'outbox_controller.dart';
import 'student_work.dart';
import 'study_repository.dart';
import 'submission.dart';
import 'submission_draft.dart';
import 'submission_outbox.dart';

final studyRepositoryProvider = Provider<StudyRepository>(
  (ref) => StudyRepository(
    ref.watch(apiClientProvider),
    ref.watch(cachedReaderProvider),
  ),
);

final submissionOutboxProvider = Provider<SubmissionOutbox>(
  (ref) => SubmissionOutbox(ref.watch(sharedPreferencesProvider)),
);

final outboxProvider =
    NotifierProvider<OutboxController, List<SubmissionDraft>>(
      OutboxController.new,
    );

final myCoursesProvider = FutureProvider.autoDispose<Loaded<List<Course>>>(
  (ref) => ref.watch(studyRepositoryProvider).courses(),
);

final courseProvider = FutureProvider.autoDispose
    .family<Loaded<CourseDetails>, String>(
      (ref, id) => ref.watch(studyRepositoryProvider).course(id),
    );

final assignmentsProvider = FutureProvider.autoDispose
    .family<Loaded<List<AssignmentItem>>, String>(
      (ref, courseId) =>
          ref.watch(studyRepositoryProvider).assignments(courseId),
    );

final assignmentProvider = FutureProvider.autoDispose
    .family<Loaded<AssignmentDetails>, String>(
      (ref, id) => ref.watch(studyRepositoryProvider).assignment(id),
    );

final assignmentWorkProvider = FutureProvider.autoDispose
    .family<List<StudentWork>, String>(
      (ref, id) => ref.watch(studyRepositoryProvider).work(id),
    );

final attemptsProvider = FutureProvider.autoDispose
    .family<List<Submission>, (String, String)>(
      (ref, key) => ref.watch(studyRepositoryProvider).attempts(key.$1, key.$2),
    );

final agendaProvider = FutureProvider.autoDispose
    .family<Loaded<Agenda>, String?>(
      (ref, date) => ref.watch(studyRepositoryProvider).agenda(date),
    );

final weekProvider = FutureProvider.autoDispose
    .family<Loaded<List<Lesson>>, (String, String)>(
      (ref, range) =>
          ref.watch(studyRepositoryProvider).schedule(range.$1, range.$2),
    );
