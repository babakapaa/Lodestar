import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'profile_providers.dart';

class ProfileHeader extends ConsumerWidget {
  const ProfileHeader({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final tokens = context.tokens;
    final textTheme = Theme.of(context).textTheme;
    final account = ref.watch(accountProvider);
    final profile = ref.watch(myProfileProvider);
    final error = account.error ?? profile.error;
    if (error != null) {
      return FormError(message: ApiErrorText.of(l10n, error));
    }
    final me = account.value;
    final mine = profile.value;
    if (me == null || mine == null) {
      return const Center(child: CircularProgressIndicator());
    }
    final name = mine.displayName.isEmpty ? me.displayName : mine.displayName;
    return LodestarCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              CircleAvatar(
                radius: 28,
                backgroundColor: tokens.palette.accentSoft,
                child: Text(
                  name.isEmpty ? '?' : name.characters.first.toUpperCase(),
                  style: textTheme.titleLarge!.copyWith(
                    color: tokens.palette.accent,
                  ),
                ),
              ),
              const SizedBox(width: Spacing.l),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(name, style: textTheme.titleLarge),
                    Text(
                      me.email,
                      style: textTheme.bodyMedium!.copyWith(
                        color: tokens.palette.muted,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          if (mine.about.isNotEmpty) ...[
            const SizedBox(height: Spacing.l),
            Text(mine.about, style: textTheme.bodyMedium),
          ],
          const SizedBox(height: Spacing.l),
          LodestarButton(
            label: l10n.profileEdit,
            icon: Icons.edit_outlined,
            variant: LodestarButtonVariant.secondary,
            onPressed: () => context.push(AppRoutes.profileEdit),
          ),
        ],
      ),
    );
  }
}
