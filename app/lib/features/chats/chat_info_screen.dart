import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../people/person_picker_dialog.dart';
import 'chat.dart';
import 'chat_details.dart';
import 'chat_edit_dialog.dart';
import 'chat_kind.dart';
import 'chat_kind_icon.dart';
import 'chat_member.dart';
import 'chat_member_tile.dart';
import 'chat_remote.dart';
import 'chat_role.dart';
import 'chats_providers.dart';

class ChatInfoScreen extends ConsumerWidget {
  const ChatInfoScreen({required this.chatId, super.key});

  final String chatId;

  Future<void> _details(
    BuildContext context,
    WidgetRef ref,
    Future<ChatDetails> Function() request,
  ) async {
    await ChatRemote.run(context, () async {
      final details = await request();
      await ref.read(chatsStoreProvider).saveChat(details.chat);
    });
    ref.invalidate(chatDetailsProvider(chatId));
  }

  Future<void> _edit(BuildContext context, WidgetRef ref, Chat chat) async {
    final edited = await ChatEditDialog.show(context, chat.title, chat.about);
    if (edited == null || !context.mounted) {
      return;
    }
    final (title, about) = edited;
    await _details(
      context,
      ref,
      () => ref.read(chatsApiProvider).update(chatId, title, about),
    );
  }

  Future<void> _mute(BuildContext context, WidgetRef ref, Chat chat) =>
      ChatRemote.run(context, () async {
        final changed = await ref
            .read(chatsApiProvider)
            .mute(chatId, muted: !chat.muted);
        await ref.read(chatsStoreProvider).saveChat(changed);
      });

  Future<void> _add(BuildContext context, WidgetRef ref) async {
    final person = await PersonPickerDialog.show(context);
    if (person == null || !context.mounted) {
      return;
    }
    await _details(
      context,
      ref,
      () => ref.read(chatsApiProvider).addMembers(chatId, [person.id]),
    );
  }

  Future<void> _remove(
    BuildContext context,
    WidgetRef ref,
    String userId,
  ) async {
    await ChatRemote.run(
      context,
      () => ref.read(chatsApiProvider).removeMember(chatId, userId),
    );
    ref.invalidate(chatDetailsProvider(chatId));
  }

  Future<void> _leave(BuildContext context, WidgetRef ref, String me) async {
    final left = await ChatRemote.run(context, () async {
      await ref.read(chatsApiProvider).removeMember(chatId, me);
      await ref.read(chatsStoreProvider).forget(chatId);
    });
    if (left && context.mounted) {
      context.go(AppRoutes.chatsPath);
    }
  }

  Future<void> _block(
    BuildContext context,
    WidgetRef ref,
    String peerId, {
    required bool blocked,
  }) async {
    final api = ref.read(chatsApiProvider);
    await ChatRemote.run(
      context,
      () => blocked ? api.unblock(peerId) : api.block(peerId),
    );
    ref.invalidate(chatBlocksProvider);
  }

  static bool _outranks(ChatRole actor, ChatRole target) =>
      actor == ChatRole.owner && target != ChatRole.owner ||
      actor == ChatRole.admin && target == ChatRole.member;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final tokens = context.tokens;
    final chat = ref.watch(chatProvider(chatId)).value;
    final details = ref.watch(chatDetailsProvider(chatId));
    final me = ref.watch(myIdProvider).value;
    if (chat == null) {
      return Scaffold(
        appBar: AppBar(),
        body: const Center(child: CircularProgressIndicator()),
      );
    }
    final direct = chat.kind == ChatKind.direct;
    final manages = !direct && chat.role.manages && !chat.archived;
    final members = details.value?.members ?? const <ChatMember>[];
    final mine = [
      for (final member in members)
        if (member.userId == me) member,
    ].firstOrNull;
    final peer = chat.peerId;
    final blocked =
        direct &&
        peer != null &&
        (ref.watch(chatBlocksProvider).value ?? const []).any(
          (user) => user.userId == peer,
        );
    return Scaffold(
      appBar: AppBar(title: Text(l10n.chatsInfo)),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              ListTile(
                contentPadding: EdgeInsets.zero,
                leading: CircleAvatar(
                  backgroundColor: tokens.palette.accentSoft,
                  child: Icon(
                    ChatKindIcon.of(chat.kind),
                    color: tokens.palette.accent,
                  ),
                ),
                title: Text(
                  chat.title,
                  style: Theme.of(context).textTheme.titleMedium,
                ),
                subtitle: chat.about.isEmpty ? null : Text(chat.about),
                trailing: manages
                    ? IconButton(
                        key: const ValueKey('chat-edit-info'),
                        tooltip: l10n.chatsEditInfo,
                        icon: const Icon(Icons.edit_outlined),
                        onPressed: () => _edit(context, ref, chat),
                      )
                    : null,
              ),
              SwitchListTile(
                key: const ValueKey('chat-mute'),
                contentPadding: EdgeInsets.zero,
                title: Text(l10n.chatsMute),
                value: chat.muted,
                onChanged: (_) => _mute(context, ref, chat),
              ),
              const SizedBox(height: Spacing.l),
              if (!direct)
                Section(
                  title: l10n.chatsMembers(chat.memberCount),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      if (details.hasError)
                        FormError(
                          message: ApiErrorText.of(l10n, details.error!),
                        )
                      else if (details.value == null)
                        const Center(child: CircularProgressIndicator())
                      else
                        for (final member in members)
                          ChatMemberTile(
                            key: ValueKey(member.userId),
                            member: member,
                            me: member.userId == me,
                            onRemove:
                                manages &&
                                    member.userId != me &&
                                    !member.synced &&
                                    _outranks(chat.role, member.role)
                                ? () => _remove(context, ref, member.userId)
                                : null,
                          ),
                      if (manages) ...[
                        const SizedBox(height: Spacing.s),
                        Align(
                          alignment: Alignment.centerLeft,
                          child: OutlinedButton.icon(
                            key: const ValueKey('chat-add-member'),
                            onPressed: () => _add(context, ref),
                            icon: const Icon(Icons.person_add_alt),
                            label: Text(l10n.chatsAddMember),
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
              const SizedBox(height: Spacing.xl),
              if (direct && peer != null)
                Align(
                  alignment: Alignment.centerLeft,
                  child: OutlinedButton.icon(
                    key: const ValueKey('chat-block'),
                    onPressed: () =>
                        _block(context, ref, peer, blocked: blocked),
                    icon: Icon(blocked ? Icons.lock_open : Icons.block),
                    label: Text(blocked ? l10n.chatsUnblock : l10n.chatsBlock),
                  ),
                ),
              if (!direct &&
                  me != null &&
                  mine != null &&
                  !mine.synced &&
                  mine.role != ChatRole.owner)
                Align(
                  alignment: Alignment.centerLeft,
                  child: OutlinedButton.icon(
                    key: const ValueKey('chat-leave'),
                    onPressed: () => _leave(context, ref, me),
                    icon: const Icon(Icons.logout),
                    label: Text(l10n.chatsLeave),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
