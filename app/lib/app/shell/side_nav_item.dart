import 'package:material_ui/material_ui.dart';

import '../../design/tokens/motion.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';

class SideNavItem extends StatelessWidget {
  const SideNavItem({
    required this.icon,
    required this.selectedIcon,
    required this.label,
    required this.selected,
    required this.onTap,
    super.key,
  });

  final IconData icon;
  final IconData selectedIcon;
  final String label;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final palette = context.tokens.palette;
    final radius = context.tokens.borderRadiusSmall;
    final color = selected ? palette.accent : palette.muted;
    final style = Theme.of(context).textTheme.labelLarge!.copyWith(
      color: selected ? palette.text : palette.muted,
      fontWeight: selected ? FontWeight.w600 : FontWeight.w500,
    );
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 1),
      child: Material(
        type: MaterialType.transparency,
        child: InkWell(
          borderRadius: radius,
          onTap: onTap,
          child: AnimatedContainer(
            duration: Motion.fast,
            height: 40,
            padding: const EdgeInsets.symmetric(horizontal: Spacing.m),
            decoration: BoxDecoration(
              color: selected ? palette.surface : null,
              borderRadius: radius,
              border: selected ? Border.all(color: palette.border) : null,
              boxShadow: selected ? palette.raised : null,
            ),
            child: Row(
              children: [
                Icon(selected ? selectedIcon : icon, size: 20, color: color),
                const SizedBox(width: Spacing.m),
                Expanded(
                  child: Text(
                    label,
                    style: style,
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
