import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import 'chat.dart';
import 'chat_kind_icon.dart';
import 'chat_preview.dart';

class ChatTile extends StatelessWidget {
  const ChatTile({required this.chat, super.key});

  final Chat chat;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final locale = Localizations.localeOf(context).toLanguageTag();
    final at = (chat.last?.createdAt ?? chat.changedAt).toLocal();
    final today = DateUtils.isSameDay(at, DateTime.now());
    final time = today
        ? DateFormat.Hm(locale).format(at)
        : DateFormat.MMMd(locale).format(at);
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: CircleAvatar(
        backgroundColor: tokens.palette.accentSoft,
        child: Icon(ChatKindIcon.of(chat.kind), color: tokens.palette.accent),
      ),
      title: Row(
        children: [
          Expanded(
            child: Text(
              chat.title,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ),
          if (chat.muted)
            Padding(
              padding: const EdgeInsets.only(left: Spacing.xs),
              child: Icon(
                Icons.notifications_off_outlined,
                size: 16,
                color: tokens.palette.muted,
              ),
            ),
        ],
      ),
      subtitle: Text(
        ChatPreview.of(context, chat.last),
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
      ),
      trailing: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Text(time, style: Theme.of(context).textTheme.labelSmall),
          const SizedBox(height: Spacing.xs),
          if (chat.unread > 0)
            Badge(
              label: Text('${chat.unread}'),
              backgroundColor: chat.muted
                  ? tokens.palette.muted
                  : tokens.palette.accent,
            ),
        ],
      ),
      onTap: () => context.push(AppRoutes.chat(chat.id)),
    );
  }
}
