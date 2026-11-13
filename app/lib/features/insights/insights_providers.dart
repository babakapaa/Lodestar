import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'course_topics.dart';
import 'insights_api.dart';

final insightsApiProvider = Provider<InsightsApi>(
  (ref) => InsightsApi(ref.watch(apiClientProvider)),
);

final courseTopicsProvider = FutureProvider.autoDispose
    .family<CourseTopics, String>(
      (ref, id) => ref.watch(insightsApiProvider).topics(id),
    );
