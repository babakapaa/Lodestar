import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'role.dart';
import 'role_tile.dart';
import 'roles_providers.dart';

class RolesScreen extends ConsumerWidget {
  const RolesScreen({super.key});

  Future<void> _move(
    BuildContext context,
    WidgetRef ref,
    List<Role> movable,
    int from,
    int to,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    if (to == from) {
      return;
    }
    try {
      await ref
          .read(rolesRepositoryProvider)
          .move(movable[from].id, movable[to].position);
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(rolesProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final roles = ref.watch(rolesProvider);
    return Scaffold(
      appBar: AppBar(title: Text(l10n.rolesTitle)),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => context.push(AppRoutes.roleNew),
        icon: const Icon(Icons.add),
        label: Text(l10n.rolesCreate),
      ),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: switch (roles) {
              AsyncData(:final value) => Builder(
                builder: (context) {
                  final movable = [
                    for (final role in value)
                      if (!role.isDefault) role,
                  ];
                  final fixed = [
                    for (final role in value)
                      if (role.isDefault) role,
                  ];
                  return ListView(
                    padding: const EdgeInsets.only(bottom: 96),
                    children: [
                      Padding(
                        padding: const EdgeInsets.all(Spacing.l),
                        child: Text(l10n.rolesHint),
                      ),
                      ReorderableListView(
                        shrinkWrap: true,
                        physics: const NeverScrollableScrollPhysics(),
                        buildDefaultDragHandles: false,
                        onReorderItem: (from, to) =>
                            _move(context, ref, movable, from, to),
                        children: [
                          for (final (index, role) in movable.indexed)
                            RoleTile(
                              key: ValueKey(role.id),
                              role: role,
                              trailing: ReorderableDragStartListener(
                                index: index,
                                child: const Icon(Icons.drag_handle),
                              ),
                            ),
                        ],
                      ),
                      for (final role in fixed)
                        RoleTile(key: ValueKey(role.id), role: role),
                    ],
                  );
                },
              ),
              AsyncError(:final error) => Padding(
                padding: const EdgeInsets.all(Spacing.l),
                child: FormError(message: ApiErrorText.of(l10n, error)),
              ),
              _ => const Center(child: CircularProgressIndicator()),
            },
          ),
        ),
      ),
    );
  }
}
