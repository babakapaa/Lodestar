import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'gpu_class.dart';
import 'gpu_node.dart';
import 'sandbox_api.dart';
import 'sandbox_quota.dart';
import 'sandbox_run.dart';
import 'sandbox_template.dart';
import 'sandbox_term.dart';
import 'topup.dart';

final sandboxApiProvider = Provider<SandboxApi>(
  (ref) => SandboxApi(ref.watch(apiClientProvider)),
);

final sandboxQuotaProvider = FutureProvider.autoDispose<SandboxQuota>(
  (ref) => ref.watch(sandboxApiProvider).quota(),
);

final sandboxRunsProvider = FutureProvider.autoDispose<List<SandboxRun>>(
  (ref) => ref.watch(sandboxApiProvider).runs(),
);

final sandboxRunProvider = FutureProvider.autoDispose
    .family<SandboxRun, String>(
      (ref, id) => ref.watch(sandboxApiProvider).run(id),
    );

final sandboxTopupsProvider = FutureProvider.autoDispose<List<Topup>>(
  (ref) => ref.watch(sandboxApiProvider).topups(),
);

final gpuClassesProvider = FutureProvider.autoDispose<List<GpuClass>>(
  (ref) => ref.watch(sandboxApiProvider).classes(),
);

final sandboxTemplatesProvider =
    FutureProvider.autoDispose<List<SandboxTemplate>>(
      (ref) => ref.watch(sandboxApiProvider).templates(),
    );

final sandboxTermsProvider = FutureProvider.autoDispose<List<SandboxTerm>>(
  (ref) => ref.watch(sandboxApiProvider).terms(),
);

final gpuNodesProvider = FutureProvider.autoDispose<List<GpuNode>>(
  (ref) => ref.watch(sandboxApiProvider).nodes(),
);

final heldRunsProvider = FutureProvider.autoDispose<List<SandboxRun>>(
  (ref) => ref.watch(sandboxApiProvider).held(),
);
