import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../../core/offline/loaded.dart';
import '../../core/offline/offline_providers.dart';
import 'course_material.dart';
import 'dialog_details.dart';
import 'navigator_api.dart';
import 'navigator_dialog.dart';

final navigatorApiProvider = Provider<NavigatorApi>(
  (ref) => NavigatorApi(
    ref.watch(apiClientProvider),
    ref.watch(cachedReaderProvider),
  ),
);

final dialogsProvider =
    FutureProvider.autoDispose<Loaded<List<NavigatorDialog>>>(
      (ref) => ref.watch(navigatorApiProvider).dialogs(),
    );

final dialogProvider = FutureProvider.autoDispose
    .family<Loaded<DialogDetails>, String>(
      (ref, id) => ref.watch(navigatorApiProvider).dialog(id),
    );

final materialsProvider = FutureProvider.autoDispose
    .family<Loaded<List<CourseMaterial>>, String>(
      (ref, courseId) => ref.watch(navigatorApiProvider).materials(courseId),
    );
