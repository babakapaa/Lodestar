import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../people/people_providers.dart';
import 'role.dart';
import 'role_dot.dart';
import 'roles_providers.dart';

class PersonRolesScreen extends ConsumerWidget {
  const PersonRolesScreen({required this.userId, super.key});

  final String userId;

  Future<void> _run(
    BuildContext context,
    WidgetRef ref,
    Future<void> Function() action,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await action();
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(personRolesProvider(userId));
  }

  Future<void> _grant(
    BuildContext context,
    WidgetRef ref,
    List<Role> available,
  ) async {
    final role = await showDialog<Role>(
      context: context,
      builder: (context) => SimpleDialog(
        title: Text(context.l10n.rolesGrant),
        children: [
          for (final role in available)
            SimpleDialogOption(
              onPressed: () => Navigator.of(context).pop(role),
              child: Row(
                children: [
                  RoleDot(color: role.color),
                  const SizedBox(width: Spacing.m),
                  Text(role.name),
                ],
              ),
            ),
        ],
      ),
    );
    if (role != null && context.mounted) {
      await _run(
        context,
        ref,
        () => ref.read(rolesRepositoryProvider).assign(role.id, userId),
      );
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final name = ref.watch(personNameProvider(userId)).value ?? '';
    final roles = ref.watch(rolesProvider);
    final assignments = ref.watch(personRolesProvider(userId));
    final error = roles.error ?? assignments.error;
    final list = roles.value;
    final held = assignments.value;
    return Scaffold(
      appBar: AppBar(title: Text(name)),
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
                : list == null || held == null
                ? const Center(child: CircularProgressIndicator())
                : Builder(
                    builder: (context) {
                      final byId = {for (final role in list) role.id: role};
                      final heldIds = {
                        for (final item in held)
                          if (item.wholeUniversity) item.roleId,
                      };
                      final available = [
                        for (final role in list)
                          if (!role.isDefault && !heldIds.contains(role.id))
                            role,
                      ];
                      return ListView(
                        padding: const EdgeInsets.all(Spacing.l),
                        children: [
                          if (held.isEmpty) Text(l10n.rolesNone),
                          for (final item in held)
                            ListTile(
                              contentPadding: EdgeInsets.zero,
                              leading: RoleDot(
                                color: byId[item.roleId]?.color ?? '',
                                size: 16,
                              ),
                              title: Text(byId[item.roleId]?.name ?? ''),
                              subtitle: item.wholeUniversity
                                  ? null
                                  : Text(l10n.rolesInScope(item.scope)),
                              trailing: IconButton(
                                tooltip: l10n.rolesRevoke,
                                icon: const Icon(Icons.close),
                                onPressed: () => _run(
                                  context,
                                  ref,
                                  () => ref
                                      .read(rolesRepositoryProvider)
                                      .unassign(item.id),
                                ),
                              ),
                            ),
                          const SizedBox(height: Spacing.l),
                          LodestarButton(
                            label: l10n.rolesGrant,
                            icon: Icons.add,
                            variant: LodestarButtonVariant.secondary,
                            onPressed: available.isEmpty
                                ? null
                                : () => _grant(context, ref, available),
                          ),
                        ],
                      );
                    },
                  ),
          ),
        ),
      ),
    );
  }
}
