import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'chat_file.dart';
import 'chat_file_link.dart';
import 'chat_message.dart';
import 'reply_quote.dart';

class ChatBubble extends StatelessWidget {
  const ChatBubble({
    required this.message,
    required this.mine,
    required this.author,
    this.quote,
    this.onLongPress,
    this.onFile,
    super.key,
  });

  final ChatMessage message;
  final bool mine;
  final String? author;
  final ReplyQuote? quote;
  final VoidCallback? onLongPress;
  final ValueChanged<ChatFile>? onFile;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final tokens = context.tokens;
    final theme = Theme.of(context);
    final small = theme.textTheme.labelSmall?.copyWith(
      color: tokens.palette.muted,
    );
    final locale = Localizations.localeOf(context).toLanguageTag();
    final views = message.views;
    return Align(
      alignment: mine
          ? AlignmentDirectional.centerEnd
          : AlignmentDirectional.centerStart,
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 520),
        child: GestureDetector(
          onLongPress: message.deleted ? null : onLongPress,
          onSecondaryTap: message.deleted ? null : onLongPress,
          child: DecoratedBox(
            decoration: BoxDecoration(
              color: mine ? tokens.palette.accentSoft : tokens.palette.surface,
              borderRadius: tokens.borderRadius,
              border: Border.all(color: tokens.palette.border),
            ),
            child: Padding(
              padding: const EdgeInsets.all(Spacing.m),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  if (author != null && !mine)
                    Text(
                      author!,
                      style: theme.textTheme.labelMedium?.copyWith(
                        color: tokens.palette.accent,
                      ),
                    ),
                  if (quote != null) ...[
                    const SizedBox(height: Spacing.xs),
                    quote!,
                    const SizedBox(height: Spacing.xs),
                  ],
                  if (message.deleted)
                    Text(
                      l10n.chatsDeleted,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        fontStyle: FontStyle.italic,
                        color: tokens.palette.muted,
                      ),
                    )
                  else if (message.text.isNotEmpty)
                    Text(message.text),
                  if (!message.deleted)
                    for (final file in message.files)
                      Padding(
                        padding: const EdgeInsets.only(top: Spacing.xs),
                        child: ChatFileLink(
                          file: file,
                          onTap: onFile == null ? null : () => onFile!(file),
                        ),
                      ),
                  const SizedBox(height: Spacing.xs),
                  Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      if (message.editedAt != null && !message.deleted) ...[
                        Text(l10n.chatsEdited, style: small),
                        const SizedBox(width: Spacing.s),
                      ],
                      Text(
                        DateFormat.Hm(locale)
                            .format(message.createdAt.toLocal()),
                        style: small,
                      ),
                      if (views != null) ...[
                        const SizedBox(width: Spacing.s),
                        Icon(
                          Icons.visibility_outlined,
                          size: 14,
                          color: tokens.palette.muted,
                        ),
                        const SizedBox(width: Spacing.xs),
                        Text('$views', style: small),
                      ],
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
