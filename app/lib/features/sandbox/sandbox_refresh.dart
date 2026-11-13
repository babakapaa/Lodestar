import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'sandbox_providers.dart';

abstract final class SandboxRefresh {
  static void all(WidgetRef ref) {
    ref
      ..invalidate(sandboxQuotaProvider)
      ..invalidate(sandboxRunsProvider)
      ..invalidate(sandboxRunProvider);
  }
}
