import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'lodestar_mark.dart';

class BrandMark extends StatelessWidget {
  const BrandMark({required this.showName, this.width, super.key});

  final bool showName;
  final double? width;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final mark = LodestarMark(size: 32, color: tokens.palette.accent);
    if (!showName) {
      return mark;
    }
    return SizedBox(
      width: width,
      child: Row(
        mainAxisSize: MainAxisSize.min,
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
