import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/tab_page.dart';
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
import 'new_chat_actions.dart';

class ChatsScreen extends ConsumerStatefulWidget {
  const ChatsScreen({super.key});

  @override
  ConsumerState<ChatsScreen> createState() => _ChatsScreenState();
}

class _ChatsScreenState extends ConsumerState<ChatsScreen> {
  ChatFolder _folder = ChatFolder.all;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final status = ref.watch(chatsSyncStatusProvider);
    final all = ref.watch(chatsProvider).value;
    final requests = [
      for (final chat in all ?? const <Chat>[])
        if (chat.request) chat,
    ];
    final folder = requests.isEmpty ? ChatFolder.all : _folder;
    final shown = [
      for (final chat in all ?? const <Chat>[])
        if (chat.request == (folder == ChatFolder.requests)) chat,
    ];
    return TabPage(
      title: l10n.tabChats,
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => NewChatActions.choose(context, ref),
        icon: const Icon(Icons.edit_outlined),
        label: Text(l10n.chatsNew),
      ),
      child: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: all == null
              ? const Center(child: CircularProgressIndicator())
              : RefreshIndicator(
                  onRefresh: ref.read(chatsSyncStatusProvider.notifier).sync,
                  child: ListView(
                    padding: const EdgeInsets.fromLTRB(
                      Spacing.l,
                      Spacing.l,
                      Spacing.l,
                      96,
                    ),
                    children: [
                      if (status == ChatsSyncStatus.offline) ...[
                        NoticeBanner(
                          icon: Icons.cloud_off_outlined,
                          message: l10n.chatsOffline,
                        ),
                        const SizedBox(height: Spacing.l),
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
                          onSelected: (value) =>
                              setState(() => _folder = value),
                        ),
                        const SizedBox(height: Spacing.s),
                      ],
                      if (shown.isEmpty)
                        EmptyState(
                          icon: Icons.chat_bubble_outline,
                          title: l10n.emptyChatsTitle,
                          body: l10n.emptyChatsBody,
                        )
                      else
                        for (final chat in shown)
                          ChatTile(key: ValueKey(chat.id), chat: chat),
                    ],
                  ),
                ),
        ),
      ),
    );
  }
}
