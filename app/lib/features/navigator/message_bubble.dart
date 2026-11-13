import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import 'citation_list.dart';
import 'navigator_labels_context.dart';
import 'navigator_message.dart';
import 'support_card.dart';

class MessageBubble extends StatelessWidget {
  const MessageBubble({required this.message, super.key});

  final NavigatorMessage message;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final step = message.step;
    if (!message.fromNavigator) {
      return Align(
        alignment: AlignmentDirectional.centerEnd,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 560),
          child: DecoratedBox(
            decoration: BoxDecoration(
              color: tokens.palette.accentSoft,
              borderRadius: tokens.borderRadius,
            ),
            child: Padding(
              padding: const EdgeInsets.all(Spacing.m),
              child: SelectableText(message.text),
            ),
          ),
        ),
      );
    }
    return LodestarCard(
      padding: const EdgeInsets.all(Spacing.m),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (step != null) ...[
            Text(
              context.navigatorLabels.step(step),
              style: Theme.of(context).textTheme.labelMedium!
                  .copyWith(color: tokens.palette.muted),
            ),
            const SizedBox(height: Spacing.s),
          ],
          SelectableText(message.text),
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
    );
  }
}
