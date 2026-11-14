import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/choice_group.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'chat.dart';
import 'chat_folder.dart';
import 'chat_tile.dart';
import 'chats_providers.dart';
import 'chats_sync_status.dart';

class ChatList extends ConsumerStatefulWidget {
  const ChatList({
    this.selected,
    this.onOpen,
    this.action,
    this.padding = const EdgeInsets.all(Spacing.m),
    super.key,
  });

  final String? selected;
  final ValueChanged<String>? onOpen;
  final Widget? action;
  final EdgeInsets padding;

  @override
  ConsumerState<ChatList> createState() => _ChatListState();
}

class _ChatListState extends ConsumerState<ChatList> {
  ChatFolder _folder = ChatFolder.all;
  String _query = '';

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final status = ref.watch(chatsSyncStatusProvider);
    final all = ref.watch(chatsProvider).value;
    if (all == null) {
      return const Center(child: CircularProgressIndicator());
    }
    final requests = [
      for (final chat in all)
        if (chat.request) chat,
    ];
    final folder = requests.isEmpty ? ChatFolder.all : _folder;
    final needle = _query.trim().toLowerCase();
    final shown = [
      for (final chat in all)
        if (chat.request == (folder == ChatFolder.requests) &&
            (needle.isEmpty || chat.title.toLowerCase().contains(needle)))
          chat,
    ];
    return RefreshIndicator(
      onRefresh: ref.read(chatsSyncStatusProvider.notifier).sync,
      child: ListView(
        padding: widget.padding,
        children: [
          if (status == ChatsSyncStatus.offline) ...[
            NoticeBanner(
              icon: Icons.cloud_off_outlined,
              message: l10n.chatsOffline,
            ),
            const SizedBox(height: Spacing.m),
          ],
          if (all.isNotEmpty || widget.action != null) ...[
            Row(
              children: [
                Expanded(
                  child: TextField(
                    decoration: InputDecoration(
                      prefixIcon: const Icon(Icons.search, size: 18),
                      hintText: l10n.chatsSearch,
                    ),
                    onChanged: (value) => setState(() => _query = value),
                  ),
                ),
                if (widget.action != null) ...[
                  const SizedBox(width: Spacing.xs),
                  widget.action!,
                ],
              ],
            ),
            const SizedBox(height: Spacing.s),
          ],
          if (requests.isNotEmpty) ...[
            ChoiceGroup<ChatFolder>(
              values: ChatFolder.values,
              selected: folder,
              labelOf: (value) => switch (value) {
                ChatFolder.all => l10n.chatsAll,
                ChatFolder.requests =>
                  '${l10n.chatsRequests} ${requests.length}',
              },
              onSelected: (value) => setState(() => _folder = value),
            ),
            const SizedBox(height: Spacing.s),
          ],
          if (all.isEmpty)
            EmptyState(
              icon: Icons.chat_bubble_outline,
              title: l10n.emptyChatsTitle,
              body: l10n.emptyChatsBody,
            )
          else
            for (final Chat chat in shown)
              ChatTile(
                key: ValueKey(chat.id),
                chat: chat,
                selected: chat.id == widget.selected,
                onOpen: widget.onOpen,
              ),
        ],
      ),
    );
  }
}
