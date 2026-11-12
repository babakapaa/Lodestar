import 'package:material_ui/material_ui.dart';

import '../tokens/motion.dart';
import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';

class LodestarChip extends StatelessWidget {
  const LodestarChip({
    required this.label,
    required this.selected,
    required this.onSelected,
    this.leading,
    super.key,
  });

  final String label;
  final bool selected;
  final VoidCallback onSelected;
  final Widget? leading;

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    final textTheme = Theme.of(context).textTheme;
    final shape = RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(999),
      side: BorderSide(color: selected ? palette.accent : palette.border),
    );
    return Semantics(
      selected: selected,
      button: true,
      child: AnimatedContainer(
        duration: Motion.of(context, Motion.fast),
        curve: Motion.curve,
        decoration: ShapeDecoration(
          color: selected ? palette.accentSoft : Colors.transparent,
          shape: shape,
        ),
        child: Material(
          type: MaterialType.transparency,
          shape: shape,
          clipBehavior: Clip.antiAlias,
          child: InkWell(
            onTap: onSelected,
            child: ConstrainedBox(
              constraints: const BoxConstraints(minHeight: 44),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: Spacing.l),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    if (leading != null) ...[
                      leading!,
                      const SizedBox(width: Spacing.s),
                    ],
                    Text(
                      label,
                      style: textTheme.labelMedium!.copyWith(
                        color: selected ? palette.accent : palette.text,
                        fontWeight: selected
                            ? FontWeight.w600
                            : FontWeight.w500,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
