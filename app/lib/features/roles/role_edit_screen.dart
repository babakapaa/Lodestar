import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'role.dart';
import 'role_form.dart';
import 'roles_providers.dart';

class RoleEditScreen extends ConsumerWidget {
  const RoleEditScreen({this.roleId, super.key});

  final String? roleId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final roles = ref.watch(rolesProvider);
    final catalog = ref.watch(permissionCatalogProvider);
    final error = roles.error ?? catalog.error;
    final list = roles.value;
    final permissions = catalog.value;
    Role? role;
    for (final item in list ?? const <Role>[]) {
      if (item.id == roleId) {
        role = item;
      }
    }
    return Scaffold(
      appBar: AppBar(
        title: Text(roleId == null ? l10n.rolesCreate : l10n.rolesEdit),
      ),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: error != null
                ? Padding(
                    padding: const EdgeInsets.all(Spacing.l),
                    child: FormError(message: ApiErrorText.of(l10n, error)),
                  )
                : list == null || permissions == null
                ? const Center(child: CircularProgressIndicator())
                : RoleForm(
                    key: ValueKey(role?.id),
                    role: role,
                    catalog: permissions,
                  ),
          ),
        ),
      ),
    );
  }
}
