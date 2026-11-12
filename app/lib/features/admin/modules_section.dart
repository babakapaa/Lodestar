import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'admin_providers.dart';

class ModulesSection extends ConsumerWidget {
  const ModulesSection({super.key});

  Future<void> _set(
    BuildContext context,
    WidgetRef ref,
    String name,
    bool enabled,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await ref.read(adminRepositoryProvider).setModule(name, enabled: enabled);
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(modulesProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final modules = ref.watch(modulesProvider);
    return Section(
      title: l10n.adminModulesTitle,
      child: switch (modules) {
        AsyncData(:final value) when value.isEmpty => Text(
          l10n.adminModulesEmpty,
        ),
        AsyncData(:final value) => Column(
          children: [
            for (final module in value)
              SwitchListTile(
                contentPadding: EdgeInsets.zero,
                title: Text(module.name),
                subtitle: module.optional
                    ? null
                    : Text(l10n.adminModuleRequired),
                value: module.enabled,
                onChanged: module.optional
                    ? (on) => _set(context, ref, module.name, on)
                    : null,
              ),
          ],
        ),
        AsyncError(:final error) => FormError(
          message: ApiErrorText.of(l10n, error),
        ),
        _ => const Center(child: CircularProgressIndicator()),
      },
    );
  }
}
