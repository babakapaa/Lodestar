import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../admin/admin_providers.dart';
import '../study/navigator_tile.dart';

class SandboxWorkTile extends ConsumerWidget {
  const SandboxWorkTile({
    required this.kind,
    required this.id,
    required this.title,
    super.key,
  });

  final String kind;
  final String id;
  final String title;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final enabled =
        ref
            .watch(modulesProvider)
            .value
            ?.any((m) => m.name == 'sandbox' && m.enabled) ??
        false;
    if (!enabled) {
      return const SizedBox.shrink();
    }
    final l10n = context.l10n;
    return Padding(
      padding: const EdgeInsets.only(bottom: Spacing.l),
      child: NavigatorTile(
        key: const ValueKey('work-sandbox'),
        icon: Icons.memory,
        title: l10n.sandboxForWorkTile,
        subtitle: l10n.sandboxForWorkHint,
        route: AppRoutes.sandboxNewFor(kind, id, title),
      ),
    );
  }
}
