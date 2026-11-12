import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/access/access_providers.dart';
import '../../core/access/permission_set.dart';
import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'unit.dart';
import 'unit_action.dart';
import 'unit_form_dialog.dart';
import 'unit_kinds.dart';
import 'unit_members_section.dart';
import 'unit_picker_dialog.dart';
import 'units_providers.dart';
import 'units_query.dart';

class UnitsScreen extends ConsumerStatefulWidget {
  const UnitsScreen({this.unitId, super.key});

  final String? unitId;

  @override
  ConsumerState<UnitsScreen> createState() => _UnitsScreenState();
}

class _UnitsScreenState extends ConsumerState<UnitsScreen> {
  bool _archived = false;

  UnitsQuery get _query =>
      UnitsQuery(parentId: widget.unitId, archived: _archived);

  Future<void> _run(Future<void> Function() action) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await action();
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(unitChildrenProvider);
    final id = widget.unitId;
    if (id != null) {
      ref.invalidate(unitDetailsProvider(id));
    }
  }

  Future<void> _create() async {
    final input = await UnitFormDialog.show(context, parentId: widget.unitId);
    if (input != null) {
      await _run(() => ref.read(unitsRepositoryProvider).create(input));
    }
  }

  Future<void> _act(UnitAction action, Unit unit) async {
    final repository = ref.read(unitsRepositoryProvider);
    switch (action) {
      case UnitAction.edit:
        final input = await UnitFormDialog.show(context, unit: unit);
        if (input != null) {
          await _run(() => repository.update(unit.id, input));
        }
      case UnitAction.move:
        final target = await UnitPickerDialog.show(context, unit);
        if (target != null) {
          await _run(() => repository.move(unit.id, target.parentId));
        }
      case UnitAction.archive:
        await _run(() => repository.archive(unit.id));
      case UnitAction.restore:
        await _run(() => repository.restore(unit.id));
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final permissions =
        ref.watch(myPermissionsProvider).value ?? PermissionSet.empty;
    final canStructure = permissions.has(PermissionSet.manageStructure);
    final canMembers = permissions.has(PermissionSet.manageMembers);
    final id = widget.unitId;
    final details = id == null ? null : ref.watch(unitDetailsProvider(id));
    final unit = details?.value?.unit;
    final path = details?.value?.path ?? const <Unit>[];
    final children = ref.watch(unitChildrenProvider(_query));
    return Scaffold(
      appBar: AppBar(
        title: Text(unit?.name ?? l10n.unitsTitle),
        actions: [
          if (canStructure && unit != null)
            PopupMenuButton<UnitAction>(
              onSelected: (action) => _act(action, unit),
              itemBuilder: (context) => [
                PopupMenuItem(
                  value: UnitAction.edit,
                  child: Text(l10n.unitsEdit),
                ),
                PopupMenuItem(
                  value: UnitAction.move,
                  child: Text(l10n.unitsMove),
                ),
                PopupMenuItem(
                  value: unit.archived
                      ? UnitAction.restore
                      : UnitAction.archive,
                  child: Text(
                    unit.archived ? l10n.unitsRestore : l10n.unitsArchive,
                  ),
                ),
              ],
            ),
        ],
      ),
      floatingActionButton: canStructure && !(unit?.archived ?? false)
          ? FloatingActionButton.extended(
              onPressed: _create,
              icon: const Icon(Icons.add),
              label: Text(l10n.unitsCreate),
            )
          : null,
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: ListView(
              padding: const EdgeInsets.fromLTRB(
                Spacing.l,
                Spacing.l,
                Spacing.l,
                96,
              ),
              children: [
                if (unit != null) ...[
                  Text(
                    [
                      for (final parent in path) parent.name,
                      UnitKinds.label(l10n, unit.kind),
                    ].join(' / '),
                    style: Theme.of(context).textTheme.bodyMedium!
                        .copyWith(color: context.tokens.palette.muted),
                  ),
                  if (unit.archived) Text(l10n.unitsArchived),
                  const SizedBox(height: Spacing.l),
                ],
                Section(
                  title: unit == null ? l10n.unitsTop : l10n.unitsInside,
                  child: Column(
                    children: [
                      switch (children) {
                        AsyncData(:final value) when value.isEmpty => Align(
                          alignment: Alignment.centerLeft,
                          child: Text(l10n.unitsEmpty),
                        ),
                        AsyncData(:final value) => Column(
                          children: [
                            for (final child in value)
                              ListTile(
                                contentPadding: EdgeInsets.zero,
                                leading: Icon(
                                  child.archived
                                      ? Icons.inventory_2_outlined
                                      : Icons.account_tree_outlined,
                                ),
                                title: Text(child.name),
                                subtitle: Text(
                                  UnitKinds.label(l10n, child.kind),
                                ),
                                trailing: const Icon(Icons.chevron_right),
                                onTap: () =>
                                    context.push(AppRoutes.unit(child.id)),
                              ),
                          ],
                        ),
                        AsyncError(:final error) => FormError(
                          message: ApiErrorText.of(l10n, error),
                        ),
                        _ => const Center(child: CircularProgressIndicator()),
                      },
                      SwitchListTile(
                        contentPadding: EdgeInsets.zero,
                        title: Text(l10n.unitsShowArchived),
                        value: _archived,
                        onChanged: (on) => setState(() => _archived = on),
                      ),
                    ],
                  ),
                ),
                if (id != null) ...[
                  const SizedBox(height: Spacing.xl),
                  UnitMembersSection(
                    unitId: id,
                    canManage: canMembers || canStructure,
                  ),
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }
}
