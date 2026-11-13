import 'package:material_ui/material_ui.dart';

import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class PendingBubble extends StatelessWidget {
  const PendingBubble({
    required this.text,
    required this.error,
    required this.onRetry,
    super.key,
  });

  final String text;
  final String? error;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final tokens = context.tokens;
    final failure = error;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (text.isNotEmpty)
          Align(
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
                  child: Text(text),
                ),
              ),
            ),
          ),
        const SizedBox(height: Spacing.m),
        if (failure == null)
          Semantics(
            liveRegion: true,
            child: Row(
              children: [
                const SizedBox.square(
                  dimension: 18,
                  child: CircularProgressIndicator(strokeWidth: 2),
                ),
                const SizedBox(width: Spacing.m),
                Expanded(
                  child: Text(
                    l10n.navigatorThinking,
                    style: TextStyle(color: tokens.palette.muted),
                  ),
                ),
              ],
            ),
          )
        else ...[
          FormError(message: failure),
          Align(
            alignment: AlignmentDirectional.centerStart,
            child: TextButton.icon(
              onPressed: onRetry,
              icon: const Icon(Icons.refresh),
              label: Text(l10n.navigatorRetry),
            ),
          ),
        ],
      ],
    );
  }
}
