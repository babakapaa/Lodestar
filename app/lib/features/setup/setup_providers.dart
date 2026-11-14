import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../../core/auth/auth_state.dart';
import 'hardware_report.dart';
import 'pilot_course.dart';
import 'setup_api.dart';
import 'setup_overview.dart';
import 'setup_status.dart';

final setupApiProvider = Provider<SetupApi>(
  (ref) => SetupApi(ref.watch(apiClientProvider)),
);

final setupStatusProvider = FutureProvider<SetupStatus>((ref) async {
  if (ref.watch(authProvider) is! SignedIn) {
    return SetupStatus.open;
  }
  return ref.watch(setupApiProvider).status();
});

final pilotCoursesProvider = FutureProvider.autoDispose<List<PilotCourse>>(
  (ref) => ref.watch(setupApiProvider).pilotCourses(),
);

final setupOverviewProvider = FutureProvider.autoDispose<SetupOverview>(
  (ref) => ref.watch(setupApiProvider).overview(),
);

final setupHardwareProvider = FutureProvider.autoDispose<HardwareReport>(
  (ref) => ref.watch(setupApiProvider).hardware(),
);
