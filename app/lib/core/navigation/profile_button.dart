import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'app_routes.dart';

class ProfileButton extends StatelessWidget {
  const ProfileButton({required this.showLabel, super.key});

  final bool showLabel;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final label = context.l10n.profileAndSettings;
    final avatar = CircleAvatar(
      radius: 16,
      backgroundColor: tokens.palette.accentSoft,
      child: Icon(Icons.person_outline, size: 18, color: tokens.palette.accent),
    );
    void open() => context.push(AppRoutes.settings);

    if (!showLabel) {
      return IconButton(tooltip: label, onPressed: open, icon: avatar);
    }
    return SizedBox(
      width: 216,
      child: Material(
        type: MaterialType.transparency,
        child: InkWell(
          borderRadius: tokens.borderRadiusSmall,
          onTap: open,
          child: Padding(
            padding: const EdgeInsets.all(Spacing.s),
            child: Row(
              children: [
                avatar,
                const SizedBox(width: Spacing.m),
                Expanded(
                  child: Text(
                    label,
                    style: Theme.of(context).textTheme.labelMedium,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
