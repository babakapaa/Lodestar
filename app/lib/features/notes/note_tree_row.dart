import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';

class NoteTreeRow extends StatelessWidget {
  const NoteTreeRow({
    required this.depth,
    required this.icon,
    required this.label,
    required this.onTap,
    this.selected = false,
    this.strong = false,
    super.key,
  });

  static const indent = 16.0;

  final int depth;
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool selected;
  final bool strong;

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    final style = Theme.of(context).textTheme.bodyMedium!.copyWith(
      color: selected ? palette.accent : palette.text,
      fontWeight: strong || selected ? FontWeight.w600 : null,
    );
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: Spacing.s),
      child: Material(
        color: selected ? palette.accentSoft : Colors.transparent,
        borderRadius: context.tokens.borderRadiusSmall,
        child: InkWell(
          borderRadius: context.tokens.borderRadiusSmall,
          onTap: onTap,
          child: SizedBox(
            height: 30,
            child: Row(
              children: [
                for (var level = 0; level < depth; level++)
                  Container(
                    width: indent,
                    margin: const EdgeInsets.only(left: Spacing.s - 1),
                    decoration: BoxDecoration(
                      border: Border(left: BorderSide(color: palette.border)),
                    ),
                  ),
                const SizedBox(width: Spacing.s),
                Icon(
                  icon,
                  size: 16,
                  color: selected ? palette.accent : palette.muted,
                ),
                const SizedBox(width: Spacing.s),
                Expanded(
                  child: Text(
                    label,
                    style: style,
                    maxLines: 1,
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
