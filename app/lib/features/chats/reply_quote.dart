import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';

class ReplyQuote extends StatelessWidget {
  const ReplyQuote({required this.author, required this.text, super.key});

  final String author;
  final String text;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final style = Theme.of(context).textTheme.bodySmall;
    return DecoratedBox(
      decoration: BoxDecoration(
        border: BorderDirectional(
          start: BorderSide(color: tokens.palette.accent, width: 3),
        ),
      ),
      child: Padding(
        padding: const EdgeInsetsDirectional.only(start: Spacing.s),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              author,
              style: style?.copyWith(color: tokens.palette.accent),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            Text(
              text,
              style: style,
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ),
      ),
    );
  }
}
