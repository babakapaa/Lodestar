import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class BrandMark extends StatelessWidget {
  const BrandMark({required this.showName, super.key});

  final bool showName;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final mark = Container(
      width: 32,
      height: 32,
      decoration: BoxDecoration(
        color: tokens.palette.accent,
        borderRadius: tokens.borderRadiusSmall,
      ),
      child: Icon(Icons.explore, color: tokens.palette.onAccent, size: 18),
    );
    if (!showName) {
      return mark;
    }
    return SizedBox(
      width: 216,
      child: Row(
        children: [
          mark,
          const SizedBox(width: Spacing.m),
          Text(
            context.l10n.appTitle,
            style: Theme.of(context).textTheme.titleLarge,
          ),
        ],
      ),
    );
  }
}
