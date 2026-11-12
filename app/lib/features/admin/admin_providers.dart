import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../../core/settings/preference_value.dart';
import 'admin_repository.dart';
import 'audit_controller.dart';
import 'audit_page.dart';
import 'module_state.dart';
import 'setting_definition.dart';

final adminRepositoryProvider = Provider<AdminRepository>(
  (ref) => AdminRepository(ref.watch(apiClientProvider)),
);

final settingsSchemaProvider =
    FutureProvider.autoDispose<List<SettingDefinition>>(
      (ref) => ref.watch(adminRepositoryProvider).schema(),
    );

final tenantSettingsProvider =
    FutureProvider.autoDispose<Map<String, PreferenceValue>>(
      (ref) => ref.watch(adminRepositoryProvider).values(),
    );

final modulesProvider = FutureProvider.autoDispose<List<ModuleState>>(
  (ref) => ref.watch(adminRepositoryProvider).modules(),
);

final auditProvider =
    AsyncNotifierProvider.autoDispose<AuditController, AuditPage>(
      AuditController.new,
    );
