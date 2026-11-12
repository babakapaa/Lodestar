import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'deletion.dart';
import 'device_export_saver.dart';
import 'export_saver.dart';
import 'privacy_repository.dart';

final privacyRepositoryProvider = Provider<PrivacyRepository>(
  (ref) => PrivacyRepository(ref.watch(apiClientProvider)),
);

final exportSaverProvider = Provider<ExportSaver>(
  (ref) => const DeviceExportSaver(),
);

final deletionProvider = FutureProvider.autoDispose<Deletion?>(
  (ref) => ref.watch(privacyRepositoryProvider).deletion(),
);
