import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'navigator_intent.dart';
import 'navigator_message.dart';

class ReplyActions extends StatelessWidget {
  const ReplyActions({
    required this.reply,
    required this.onMore,
    required this.onMisread,
    super.key,
  });

  final NavigatorMessage reply;
  final VoidCallback onMore;
  final VoidCallback onMisread;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    if (reply.intent == NavigatorIntent.closed) {
      return const SizedBox.shrink();
    }
    return Wrap(
      spacing: Spacing.s,
      crossAxisAlignment: WrapCrossAlignment.center,
      children: [
        if (reply.step != null && !reply.lastStep)
          TextButton.icon(
            onPressed: onMore,
            icon: const Icon(Icons.stairs_outlined),
            label: Text(l10n.navigatorMore),
          ),
        TextButton.icon(
          onPressed: onMisread,
          icon: const Icon(Icons.help_outline),
          label: Text(l10n.navigatorMisread),
        ),
        if (reply.lastStep)
          Text(
            l10n.navigatorLastStep,
            style: TextStyle(color: context.palette.muted),
          ),
      ],
    );
  }
}
