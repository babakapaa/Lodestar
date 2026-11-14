import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';

class PaneHeading extends StatelessWidget {
  const PaneHeading({required this.title, this.trailing, super.key});

  final String title;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    final style = Theme.of(context).textTheme.labelMedium!.copyWith(
      color: context.palette.muted,
      fontWeight: FontWeight.w600,
      letterSpacing: 0.6,
    );
    return Padding(
      padding: const EdgeInsets.fromLTRB(
        Spacing.l,
        Spacing.l,
        Spacing.s,
        Spacing.xs,
      ),
      child: Row(
        children: [
          Expanded(child: Text(title.toUpperCase(), style: style)),
          ?trailing,
        ],
      ),
    );
  }
}
