import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import 'chat.dart';
import 'chat_avatar.dart';
import 'chat_kind.dart';
import 'chat_kind_icon.dart';
import 'chat_preview.dart';

class ChatTile extends StatelessWidget {
  const ChatTile({
    required this.chat,
    this.selected = false,
    this.onOpen,
    super.key,
  });

  final Chat chat;
  final bool selected;
  final ValueChanged<String>? onOpen;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    final textTheme = Theme.of(context).textTheme;
    final locale = Localizations.localeOf(context).toLanguageTag();
    final at = (chat.last?.createdAt ?? chat.changedAt).toLocal();
    final today = DateUtils.isSameDay(at, DateTime.now());
    final time = today
        ? DateFormat.Hm(locale).format(at)
        : DateFormat.MMMd(locale).format(at);
    final muted = textTheme.bodyMedium!.copyWith(color: palette.muted);
    return Material(
      color: selected ? palette.accentSoft : Colors.transparent,
      borderRadius: tokens.borderRadiusSmall,
      child: InkWell(
        borderRadius: tokens.borderRadiusSmall,
        onTap: () => onOpen == null
            ? context.push(AppRoutes.chat(chat.id))
            : onOpen!(chat.id),
        child: Padding(
          padding: const EdgeInsets.symmetric(
            horizontal: Spacing.s,
            vertical: Spacing.s,
          ),
          child: Row(
            children: [
              ChatAvatar(title: chat.title),
              const SizedBox(width: Spacing.m),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        if (chat.kind != ChatKind.direct) ...[
                          Icon(
                            ChatKindIcon.of(chat.kind),
                            size: 16,
                            color: palette.muted,
                          ),
                          const SizedBox(width: Spacing.xs),
                        ],
                        Expanded(
                          child: Text(
                            chat.title,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: textTheme.titleSmall,
                          ),
                        ),
                        if (chat.muted)
                          Padding(
                            padding: const EdgeInsets.only(left: Spacing.xs),
                            child: Icon(
                              Icons.notifications_off_outlined,
                              size: 14,
                              color: palette.muted,
                            ),
                          ),
                        const SizedBox(width: Spacing.s),
                        Text(
                          time,
                          style: textTheme.labelSmall!.copyWith(
                            color: palette.muted,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 2),
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            ChatPreview.of(context, chat.last),
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: muted,
                          ),
                        ),
                        if (chat.unread > 0)
                          Badge(
                            label: Text('${chat.unread}'),
                            backgroundColor: chat.muted
                                ? palette.muted
                                : palette.accent,
                            textColor: palette.onAccent,
                          ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
