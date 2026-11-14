import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/layout_size.dart';
import '../../core/navigation/tab_page.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'chat_list.dart';
import 'chats_workspace.dart';
import 'new_chat_actions.dart';

class ChatsScreen extends ConsumerWidget {
  const ChatsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    if (LayoutSize.of(MediaQuery.sizeOf(context).width) ==
        LayoutSize.expanded) {
      return TabPage(title: l10n.tabChats, child: const ChatsWorkspace());
    }
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
          child: const ChatList(
            padding: EdgeInsets.fromLTRB(Spacing.l, Spacing.l, Spacing.l, 96),
          ),
        ),
      ),
    );
  }
}
