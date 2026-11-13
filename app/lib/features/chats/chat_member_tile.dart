import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'chat_member.dart';
import 'chat_role.dart';

class ChatMemberTile extends StatelessWidget {
  const ChatMemberTile({
    required this.member,
    required this.me,
    this.onRemove,
    super.key,
  });

  final ChatMember member;
  final bool me;
  final VoidCallback? onRemove;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final role = switch (member.role) {
      ChatRole.owner => l10n.chatsRoleOwner,
      ChatRole.admin => l10n.chatsRoleAdmin,
      ChatRole.member => null,
    };
    final notes = [?role, if (member.synced) l10n.chatsSynced];
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: const Icon(Icons.person_outline),
      title: Text(
        me ? '${member.displayName} (${l10n.chatsYou})' : member.displayName,
      ),
      subtitle: notes.isEmpty ? null : Text(notes.join(', ')),
      trailing: onRemove == null
          ? null
          : IconButton(
              key: ValueKey('remove-${member.userId}'),
              tooltip: l10n.chatsRemoveMember,
              icon: const Icon(Icons.person_remove_outlined),
              onPressed: onRemove,
            ),
    );
  }
}
