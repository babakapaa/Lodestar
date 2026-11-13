import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'assignment_item.dart';
import 'course_person.dart';
import 'course_preview_api.dart';

typedef PreviewQuery = ({String courseId, String? studentId, DateTime at});

final coursePreviewApiProvider = Provider<CoursePreviewApi>(
  (ref) => CoursePreviewApi(ref.watch(apiClientProvider)),
);

final courseStudentsProvider = FutureProvider.autoDispose
    .family<List<CoursePerson>, String>(
      (ref, courseId) => ref.watch(coursePreviewApiProvider).students(courseId),
    );

final coursePreviewProvider = FutureProvider.autoDispose
    .family<List<AssignmentItem>, PreviewQuery>(
      (ref, q) => ref
          .watch(coursePreviewApiProvider)
          .preview(q.courseId, at: q.at, studentId: q.studentId),
    );
