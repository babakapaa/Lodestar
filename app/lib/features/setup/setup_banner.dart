import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'setup_providers.dart';

class SetupBanner extends ConsumerWidget {
  const SetupBanner({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final status = ref.watch(setupStatusProvider).value;
    if (status == null || !status.takesPart) {
      return const SizedBox.shrink();
    }
    final l10n = context.l10n;
    return Padding(
      padding: const EdgeInsets.only(bottom: Spacing.xl),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          NoticeBanner(
            icon: Icons.construction_outlined,
            message: l10n.setupBanner,
          ),
          TextButton(
            key: const ValueKey('setup-banner-open'),
            onPressed: () => context.push(
              status.manages
                  ? AppRoutes.setup
                  : AppRoutes.setupSection(status.assigned.first),
            ),
            child: Text(l10n.setupOpen),
          ),
        ],
      ),
    );
  }
}
