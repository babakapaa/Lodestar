import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/empty_state.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'chat.dart';
import 'chat_kind.dart';
import 'chat_origin.dart';
import 'chat_tile.dart';
import 'chats_providers.dart';
import 'new_chat.dart';
import 'new_chat_actions.dart';

class CourseChatsScreen extends ConsumerWidget {
  const CourseChatsScreen({
    required this.courseId,
    required this.teaches,
    super.key,
  });

  final String courseId;
  final bool teaches;

  Future<void> _createGroup(BuildContext context, WidgetRef ref) =>
      NewChatActions.open(
        context,
        ref,
        () => ref
            .read(chatsApiProvider)
            .create(
              NewChat(
                kind: ChatKind.group,
                origin: ChatOrigin(kind: ChatOrigin.course, id: courseId),
              ),
            ),
      );

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    ref.watch(chatsSyncStatusProvider);
    final all = ref.watch(chatsProvider).value;
    final chats = [
      for (final chat in all ?? const <Chat>[])
        if (chat.origin?.kind == ChatOrigin.course &&
            chat.origin?.id == courseId)
          chat,
    ];
    final grouped = chats.any((chat) => chat.kind == ChatKind.group);
    return Scaffold(
      appBar: AppBar(title: Text(l10n.chatsCourseTitle)),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: all == null
              ? const Center(child: CircularProgressIndicator())
              : ListView(
                  padding: const EdgeInsets.all(Spacing.l),
                  children: [
                    if (chats.isEmpty)
                      EmptyState(
                        icon: Icons.forum_outlined,
                        title: l10n.chatsCourseEmptyTitle,
                        body: l10n.chatsCourseEmptyBody,
                      )
                    else
                      for (final chat in chats)
                        ChatTile(key: ValueKey(chat.id), chat: chat),
                    if (teaches && !grouped) ...[
                      const SizedBox(height: Spacing.l),
                      OutlinedButton.icon(
                        key: const ValueKey('course-group'),
                        onPressed: () => _createGroup(context, ref),
                        icon: const Icon(Icons.group_add_outlined),
                        label: Text(l10n.chatsCourseGroup),
                      ),
                    ],
                  ],
                ),
        ),
      ),
    );
  }
}
