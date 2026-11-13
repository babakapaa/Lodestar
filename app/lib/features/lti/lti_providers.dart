import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'browser_lms_return.dart';
import 'lms_return.dart';
import 'lti_api.dart';
import 'lti_context.dart';
import 'lti_deep_link.dart';
import 'lti_platform.dart';
import 'lti_tool.dart';

final ltiApiProvider = Provider<LtiApi>(
  (ref) => LtiApi(ref.watch(apiClientProvider)),
);

final ltiToolProvider = FutureProvider.autoDispose<LtiTool>(
  (ref) => ref.watch(ltiApiProvider).tool(),
);

final ltiPlatformsProvider = FutureProvider.autoDispose<List<LtiPlatform>>(
  (ref) => ref.watch(ltiApiProvider).platforms(),
);

final ltiContextProvider = FutureProvider.autoDispose
    .family<LtiContext, String>(
      (ref, id) => ref.watch(ltiApiProvider).context(id),
    );

final courseLmsProvider = FutureProvider.autoDispose
    .family<List<LtiContext>, String>(
      (ref, courseId) =>
          ref.watch(ltiApiProvider).byOrigin('study.course', courseId),
    );

final ltiDeepLinkProvider = FutureProvider.autoDispose
    .family<LtiDeepLink, String>(
      (ref, id) => ref.watch(ltiApiProvider).deepLink(id),
    );

final lmsReturnProvider = Provider<LmsReturn>(
  (ref) => const BrowserLmsReturn(),
);
