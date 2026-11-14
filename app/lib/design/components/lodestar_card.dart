import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';

class LodestarCard extends StatelessWidget {
  const LodestarCard({
    required this.child,
    this.padding = const EdgeInsets.all(Spacing.l),
    super.key,
  });

  static const tilePadding = EdgeInsets.symmetric(
    horizontal: Spacing.l,
    vertical: Spacing.xs,
  );

  final Widget child;
  final EdgeInsetsGeometry padding;

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
          child: Padding(padding: padding, child: child),
        ),
      ),
    );
  }
}
