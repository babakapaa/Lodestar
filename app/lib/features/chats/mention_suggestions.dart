import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';
import 'chat_member.dart';

class MentionSuggestions extends StatelessWidget {
  const MentionSuggestions({
    required this.members,
    required this.onPick,
    super.key,
  });

  final List<ChatMember> members;
  final ValueChanged<ChatMember> onPick;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    return ConstrainedBox(
      constraints: const BoxConstraints(maxHeight: 200),
      child: Material(
        color: tokens.palette.surface,
        clipBehavior: Clip.antiAlias,
        shape: RoundedRectangleBorder(
          borderRadius: tokens.borderRadius,
          side: BorderSide(color: tokens.palette.border),
        ),
        child: ListView(
          shrinkWrap: true,
          children: [
            for (final member in members)
              ListTile(
                key: ValueKey('mention-${member.userId}'),
                dense: true,
                leading: const Icon(Icons.alternate_email),
                title: Text(member.displayName),
                onTap: () => onPick(member),
              ),
          ],
        ),
      ),
    );
  }
}
