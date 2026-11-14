import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/empty_state.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'chat_list.dart';
import 'chat_screen.dart';
import 'new_chat_actions.dart';

class ChatsWorkspace extends ConsumerStatefulWidget {
  const ChatsWorkspace({super.key});

  static const listWidth = 360.0;

  @override
  ConsumerState<ChatsWorkspace> createState() => _ChatsWorkspaceState();
}

class _ChatsWorkspaceState extends ConsumerState<ChatsWorkspace> {
  String? _selected;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = context.palette;
    final selected = _selected;
    return Row(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Container(
          width: ChatsWorkspace.listWidth,
          color: palette.surface,
          child: ChatList(
            selected: selected,
            onOpen: (id) => setState(() => _selected = id),
            action: IconButton.filledTonal(
              tooltip: l10n.chatsNew,
              icon: const Icon(Icons.edit_outlined),
              onPressed: () => NewChatActions.choose(context, ref),
            ),
          ),
        ),
        const VerticalDivider(width: 1),
        Expanded(
          child: selected == null
              ? ColoredBox(
                  color: palette.side,
                  child: EmptyState(
                    icon: Icons.forum_outlined,
                    title: l10n.chatsPickTitle,
                    body: l10n.chatsPickBody,
                  ),
                )
              : ChatScreen(key: ValueKey(selected), chatId: selected),
        ),
      ],
    );
  }
}
