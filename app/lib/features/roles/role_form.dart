import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'permission_labels.dart';
import 'role.dart';
import 'role_color.dart';
import 'role_dot.dart';
import 'role_input.dart';
import 'roles_providers.dart';

class RoleForm extends ConsumerStatefulWidget {
  const RoleForm({required this.role, required this.catalog, super.key});

  final Role? role;
  final List<String> catalog;

  @override
  ConsumerState<RoleForm> createState() => _RoleFormState();
}

class _RoleFormState extends ConsumerState<RoleForm> {
  late final _name = TextEditingController(text: widget.role?.name ?? '');
  late String _color = widget.role?.color ?? '';
  late final Set<String> _permissions = {...?widget.role?.permissions};
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  Future<void> _run(Future<void> Function() action) async {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await action();
      ref.invalidate(rolesProvider);
      router.pop();
    } on ApiException catch (error) {
      if (mounted) {
        setState(() => _error = ApiErrorText.of(l10n, error));
      }
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  Future<void> _save() {
    final repository = ref.read(rolesRepositoryProvider);
    final input = RoleInput(
      name: _name.text.trim(),
      color: _color,
      permissions: _permissions,
    );
    final role = widget.role;
    return _run(
      () => role == null
          ? repository.create(input)
          : repository.update(role.id, input),
    );
  }

  Future<void> _delete(Role role) async {
    final l10n = context.l10n;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.rolesDeleteConfirm(role.name)),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: Text(l10n.actionCancel),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: Text(l10n.rolesDelete),
          ),
        ],
      ),
    );
    if (confirmed ?? false) {
      await _run(() => ref.read(rolesRepositoryProvider).delete(role.id));
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = PermissionLabels(l10n);
    final role = widget.role;
    final error = _error;
    final isDefault = role?.isDefault ?? false;
    return ListView(
      padding: const EdgeInsets.all(Spacing.l),
      children: [
        TextField(
          controller: _name,
          maxLength: 100,
          enabled: !isDefault,
          decoration: InputDecoration(labelText: l10n.rolesName),
        ),
        const SizedBox(height: Spacing.m),
        Section(
          title: l10n.rolesColor,
          child: Wrap(
            spacing: Spacing.s,
            children: [
              for (final choice in RoleColor.choices)
                ChoiceChip(
                  label: RoleDot(color: choice, size: 16),
                  tooltip: choice.isEmpty ? l10n.rolesNoColor : choice,
                  selected: _color == choice,
                  onSelected: (_) => setState(() => _color = choice),
                ),
            ],
          ),
        ),
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.rolesPermissions,
          child: Column(
            children: [
              for (final permission in widget.catalog)
                CheckboxListTile(
                  contentPadding: EdgeInsets.zero,
                  title: Text(labels.of(permission)),
                  value: _permissions.contains(permission),
                  onChanged: (on) => setState(
                    () => on ?? false
                        ? _permissions.add(permission)
                        : _permissions.remove(permission),
                  ),
                ),
            ],
          ),
        ),
        if (error != null) FormError(message: error),
        const SizedBox(height: Spacing.xl),
        LodestarButton(label: l10n.actionSave, onPressed: _busy ? null : _save),
        if (role != null && !isDefault) ...[
          const SizedBox(height: Spacing.m),
          LodestarButton(
            label: l10n.rolesDelete,
            icon: Icons.delete_outline,
            variant: LodestarButtonVariant.secondary,
            onPressed: _busy ? null : () => _delete(role),
          ),
        ],
      ],
    );
  }
}
