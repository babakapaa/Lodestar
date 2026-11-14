import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../design/components/tile_group.dart';
import '../../l10n/l10n_context.dart';
import '../chats/chat.dart';
import '../chats/chat_tile.dart';
import '../chats/chats_providers.dart';
import 'today_placeholder.dart';

class TodayMessages extends ConsumerWidget {
  const TodayMessages({super.key});

  static const limit = 5;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    ref.watch(chatsSyncStatusProvider);
    final chats = ref.watch(chatsProvider).value ?? const <Chat>[];
    final unread = chats.where((chat) => chat.unread > 0).take(limit).toList();
    return Section(
      title: l10n.todayMessages,
      child: unread.isEmpty
          ? TodayPlaceholder(
              icon: Icons.mark_chat_read_outlined,
              text: l10n.todayNoMessages,
            )
          : TileGroup(
              children: [
                for (final chat in unread)
                  ChatTile(key: ValueKey('today-chat:${chat.id}'), chat: chat),
              ],
            ),
    );
  }
}
