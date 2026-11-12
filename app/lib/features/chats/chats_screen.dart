import 'package:material_ui/material_ui.dart';

import '../../core/navigation/tab_page.dart';
import '../../design/components/empty_state.dart';
import '../../l10n/l10n_context.dart';

class ChatsScreen extends StatelessWidget {
  const ChatsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return TabPage(
      title: l10n.tabChats,
      child: EmptyState(
        icon: Icons.chat_bubble_outline,
        title: l10n.emptyChatsTitle,
        body: l10n.emptyChatsBody,
      ),
    );
  }
}
