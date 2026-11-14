import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';

class TileGroup extends StatelessWidget {
  const TileGroup({required this.children, super.key});

  final List<Widget> children;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    return DecoratedBox(
      decoration: BoxDecoration(
        color: tokens.palette.surface,
        borderRadius: tokens.borderRadius,
        border: Border.all(color: tokens.palette.border),
        boxShadow: tokens.palette.shadow,
      ),
      child: ClipRRect(
        borderRadius: tokens.borderRadius,
        child: Material(
          type: MaterialType.transparency,
          child: Column(
            children: [
              for (final (index, child) in children.indexed) ...[
                if (index > 0) const Divider(indent: Spacing.l),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: Spacing.l),
                  child: child,
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
