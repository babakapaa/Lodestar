import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'message_action.dart';

class MessageActionsSheet extends StatelessWidget {
  const MessageActionsSheet({required this.actions, super.key});

  final List<MessageAction> actions;

  static Future<MessageAction?> show(
    BuildContext context,
    List<MessageAction> actions,
  ) => showModalBottomSheet<MessageAction>(
    context: context,
    builder: (context) => MessageActionsSheet(actions: actions),
  );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return SafeArea(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          for (final action in actions)
            ListTile(
              leading: Icon(switch (action) {
                MessageAction.reply => Icons.reply,
                MessageAction.edit => Icons.edit_outlined,
                MessageAction.delete => Icons.delete_outline,
                MessageAction.copy => Icons.copy,
              }),
              title: Text(switch (action) {
                MessageAction.reply => l10n.chatsReply,
                MessageAction.edit => l10n.chatsEdit,
                MessageAction.delete => l10n.chatsDelete,
                MessageAction.copy => l10n.chatsCopy,
              }),
              onTap: () => Navigator.of(context).pop(action),
            ),
        ],
      ),
    );
  }
}
