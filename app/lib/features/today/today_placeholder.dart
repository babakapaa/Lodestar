import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';

class TodayPlaceholder extends StatelessWidget {
  const TodayPlaceholder({required this.icon, required this.text, super.key});

  final IconData icon;
  final String text;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: tokens.borderRadius,
        border: Border.all(color: palette.border),
      ),
      child: Padding(
        padding: const EdgeInsets.all(Spacing.l),
        child: Row(
          children: [
            Icon(icon, color: palette.muted, size: 20),
            const SizedBox(width: Spacing.m),
            Expanded(
              child: Text(text, style: TextStyle(color: palette.muted)),
            ),
          ],
        ),
      ),
    );
  }
}
