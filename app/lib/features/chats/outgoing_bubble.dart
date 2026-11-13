import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'chat_file_link.dart';
import 'outgoing.dart';

class OutgoingBubble extends StatelessWidget {
  const OutgoingBubble({
    required this.outgoing,
    required this.onRetry,
    required this.onDiscard,
    super.key,
  });

  final Outgoing outgoing;
  final VoidCallback onRetry;
  final VoidCallback onDiscard;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final tokens = context.tokens;
    final failed = outgoing.error != null;
    return Align(
      alignment: AlignmentDirectional.centerEnd,
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 520),
        child: DecoratedBox(
          decoration: BoxDecoration(
            color: tokens.palette.accentSoft,
            borderRadius: tokens.borderRadius,
            border: Border.all(
              color: failed ? tokens.palette.warn : tokens.palette.border,
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.all(Spacing.m),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              mainAxisSize: MainAxisSize.min,
              children: [
                if (outgoing.text.isNotEmpty) Text(outgoing.text),
                for (final file in outgoing.files)
                  Padding(
                    padding: const EdgeInsets.only(top: Spacing.xs),
                    child: ChatFileLink(file: file),
                  ),
                const SizedBox(height: Spacing.xs),
                if (failed)
                  Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        l10n.chatsFailed,
                        style: TextStyle(color: tokens.palette.warn),
                      ),
                      TextButton(
                        onPressed: onRetry,
                        child: Text(l10n.actionRetry),
                      ),
                      TextButton(
                        onPressed: onDiscard,
                        child: Text(l10n.chatsDiscard),
                      ),
                    ],
                  )
                else
                  Icon(Icons.schedule, size: 14, color: tokens.palette.muted),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
