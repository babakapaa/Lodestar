import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'chat_report.dart';

class ChatReportCard extends StatelessWidget {
  const ChatReportCard({
    required this.report,
    required this.onResolve,
    super.key,
  });

  final ChatReport report;
  final VoidCallback onResolve;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final tokens = context.tokens;
    final theme = Theme.of(context);
    final locale = Localizations.localeOf(context).toLanguageTag();
    final time = DateFormat.yMMMd(locale).add_Hm();
    return LodestarCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(report.chatTitle, style: theme.textTheme.titleMedium),
          Text(
            time.format(report.createdAt.toLocal()),
            style: theme.textTheme.labelSmall?.copyWith(
              color: tokens.palette.muted,
            ),
          ),
          const SizedBox(height: Spacing.s),
          Text(l10n.chatReportsReason(report.reason)),
          const SizedBox(height: Spacing.m),
          for (final message in report.messages)
            Padding(
              padding: const EdgeInsets.only(bottom: Spacing.s),
              child: DecoratedBox(
                decoration: BoxDecoration(
                  color: tokens.palette.surface2,
                  borderRadius: tokens.borderRadius,
                ),
                child: Padding(
                  padding: const EdgeInsets.all(Spacing.m),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        message.authorName.isEmpty
                            ? l10n.chatsUnknown
                            : message.authorName,
                        style: theme.textTheme.labelMedium?.copyWith(
                          color: tokens.palette.accent,
                        ),
                      ),
                      Text(message.text),
                    ],
                  ),
                ),
              ),
            ),
          Align(
            alignment: Alignment.centerRight,
            child: FilledButton.tonal(
              key: ValueKey('resolve-${report.id}'),
              onPressed: onResolve,
              child: Text(l10n.chatReportsResolve),
            ),
          ),
        ],
      ),
    );
  }
}
