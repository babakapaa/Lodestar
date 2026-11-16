import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../features/profile/profile_providers.dart';
import '../../l10n/l10n_context.dart';

class SideAccount extends ConsumerWidget {
  const SideAccount({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    final textTheme = Theme.of(context).textTheme;
    final label = context.l10n.profileAndSettings;
    final name = ref.watch(accountProvider).value?.displayName ?? '';
    return Material(
      type: MaterialType.transparency,
      child: InkWell(
        borderRadius: tokens.borderRadiusSmall,
        onTap: () => context.push(AppRoutes.profile),
        child: Padding(
          padding: const EdgeInsets.all(Spacing.s),
          child: Row(
            children: [
              CircleAvatar(
                radius: 17,
                backgroundColor: palette.accent,
                child: name.isEmpty
                    ? Icon(
                        Icons.person_outline,
                        size: 18,
                        color: palette.onAccent,
                      )
                    : Text(
                        name.characters.first.toUpperCase(),
                        style: textTheme.titleSmall!.copyWith(
                          color: palette.onAccent,
                        ),
                      ),
              ),
              const SizedBox(width: Spacing.m),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    if (name.isNotEmpty)
                      Text(
                        name,
                        style: textTheme.titleSmall,
                        overflow: TextOverflow.ellipsis,
                      ),
                    Text(
                      label,
                      style: name.isEmpty
                          ? textTheme.labelMedium
                          : textTheme.bodySmall,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              Icon(Icons.unfold_more, size: 18, color: palette.muted),
            ],
          ),
        ),
      ),
    );
  }
}
