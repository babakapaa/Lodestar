import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../design/markdown/markdown_view.dart';
import 'citation_list.dart';
import 'navigator_avatar.dart';
import 'navigator_labels_context.dart';
import 'navigator_message.dart';
import 'support_card.dart';

class MessageBubble extends StatelessWidget {
  const MessageBubble({required this.message, super.key});

  final NavigatorMessage message;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final textTheme = Theme.of(context).textTheme;
    final body = textTheme.bodyLarge!.copyWith(height: 1.6);
    final step = message.step;
    if (!message.fromNavigator) {
      return Align(
        alignment: AlignmentDirectional.centerEnd,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 560),
          child: DecoratedBox(
            decoration: BoxDecoration(
              color: tokens.palette.surface2,
              borderRadius: BorderRadius.circular(18),
            ),
            child: Padding(
              padding: const EdgeInsets.symmetric(
                horizontal: Spacing.l,
                vertical: Spacing.m,
              ),
              child: SelectableText(message.text, style: body),
            ),
          ),
        ),
      );
    }
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const NavigatorAvatar(),
        const SizedBox(width: Spacing.m),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (step != null) ...[
                Text(
                  context.navigatorLabels.step(step),
                  style: textTheme.labelMedium!.copyWith(
                    color: tokens.palette.accent,
                  ),
                ),
                const SizedBox(height: Spacing.xs),
              ],
              SelectionArea(child: MarkdownView(message.text, style: body)),
              if (message.support.isNotEmpty) ...[
                const SizedBox(height: Spacing.m),
                SupportCard(contacts: message.support),
              ],
              if (message.citations.isNotEmpty) ...[
                const SizedBox(height: Spacing.m),
                CitationList(citations: message.citations),
              ],
            ],
          ),
        ),
      ],
    );
  }
}
