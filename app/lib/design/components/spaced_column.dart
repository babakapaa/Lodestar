import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';

class SpacedColumn extends StatelessWidget {
  const SpacedColumn({
    required this.children,
    this.gap = Spacing.xl,
    super.key,
  });

  final List<Widget> children;
  final double gap;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        for (final (i, child) in children.indexed) ...[
          if (i > 0) SizedBox(height: gap),
          child,
        ],
      ],
    );
  }
}
