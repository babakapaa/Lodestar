import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../l10n/l10n_context.dart';
import 'role.dart';
import 'role_dot.dart';

class RoleTile extends StatelessWidget {
  const RoleTile({required this.role, this.trailing, super.key});

  final Role role;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return ListTile(
      leading: RoleDot(color: role.color, size: 16),
      title: Text(role.name),
      subtitle: Text(
        role.isDefault
            ? l10n.rolesDefault
            : l10n.rolesPermissionCount(role.permissions.length),
      ),
      trailing: trailing,
      onTap: () => context.push(AppRoutes.role(role.id)),
    );
  }
}
