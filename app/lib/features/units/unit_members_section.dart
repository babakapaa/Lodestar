import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../people/person_picker_dialog.dart';
import 'units_providers.dart';

class UnitMembersSection extends ConsumerWidget {
  const UnitMembersSection({
    required this.unitId,
    required this.canManage,
    super.key,
  });

  final String unitId;
  final bool canManage;

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
    ref.invalidate(unitMembersProvider(unitId));
  }

  Future<void> _add(BuildContext context, WidgetRef ref) async {
    final person = await PersonPickerDialog.show(context);
    if (person != null && context.mounted) {
      await _run(
        context,
        ref,
        () => ref.read(unitsRepositoryProvider).addMember(unitId, person.id),
      );
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final members = ref.watch(unitMembersProvider(unitId));
    return Section(
      title: l10n.unitsMembers,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          switch (members) {
            AsyncData(:final value) when value.isEmpty => Text(
              l10n.unitsNoMembers,
            ),
            AsyncData(:final value) => Column(
              children: [
                for (final member in value)
                  ListTile(
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.person_outline),
                    title: Text(member.displayName),
                    subtitle: Text(member.email),
                    trailing: canManage
                        ? IconButton(
                            tooltip: l10n.unitsRemoveMember,
                            icon: const Icon(Icons.close),
                            onPressed: () => _run(
                              context,
                              ref,
                              () => ref
                                  .read(unitsRepositoryProvider)
                                  .removeMember(unitId, member.userId),
                            ),
                          )
                        : null,
                  ),
              ],
            ),
            AsyncError(:final error) => FormError(
              message: ApiErrorText.of(l10n, error),
            ),
            _ => const Center(child: CircularProgressIndicator()),
          },
          if (canManage) ...[
            const SizedBox(height: Spacing.m),
            Align(
              alignment: Alignment.centerLeft,
              child: LodestarButton(
                label: l10n.unitsAddMember,
                icon: Icons.person_add_alt,
                variant: LodestarButtonVariant.secondary,
                onPressed: () => _add(context, ref),
              ),
            ),
          ],
        ],
      ),
    );
  }
}
