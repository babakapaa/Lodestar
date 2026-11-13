import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'chat.dart';
import 'chat_bubble.dart';
import 'chat_composer.dart';
import 'chat_draft.dart';
import 'chat_entry.dart';
import 'chat_file.dart';
import 'chat_files.dart';
import 'chat_kind.dart';
import 'chat_member.dart';
import 'chat_message.dart';
import 'chats_providers.dart';
import 'edit_message_dialog.dart';
import 'message_action.dart';
import 'message_actions_sheet.dart';
import 'outgoing_bubble.dart';
import 'reply_quote.dart';
import 'report_dialog.dart';

class ChatScreen extends ConsumerStatefulWidget {
  const ChatScreen({required this.chatId, super.key});

  final String chatId;

  @override
  ConsumerState<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends ConsumerState<ChatScreen> {
  ChatMessage? _replyTo;
  bool _loadingOlder = false;

  String get _id => widget.chatId;

  @override
  void initState() {
    super.initState();
    ref.listenManual(chatProvider(_id), (_, next) {
      final chat = next.value;
      if (chat != null && chat.unread > 0) {
        ref.read(chatsSyncProvider).read(chat);
      }
    }, fireImmediately: true);
  }

  Future<void> _older() async {
    setState(() => _loadingOlder = true);
    try {
      await ref.read(chatsSyncProvider).older(_id);
    } on Object catch (error) {
      _fail(error);
    } finally {
      if (mounted) {
        setState(() => _loadingOlder = false);
      }
    }
  }

  Future<void> _send(ChatDraft draft) async {
    final reply = _replyTo?.seq;
    setState(() => _replyTo = null);
    await ref
        .read(chatsStoreProvider)
        .queue(
          _id,
          draft.text,
          replyTo: reply,
          mentions: draft.mentions,
          files: draft.files,
        );
  }

  Future<List<ChatFile>> _attach() async {
    try {
      return await ChatFiles(ref).pick(ChatComposer.maxFiles);
    } on Object catch (error) {
      _fail(error);
      return const [];
    }
  }

  Future<void> _save(ChatFile file) async {
    try {
      await ChatFiles(ref).save(file);
    } on Object catch (error) {
      _fail(error);
    }
  }

  Future<void> _report(ChatMessage message) async {
    final reason = await ReportDialog.show(context);
    if (reason == null || !mounted) {
      return;
    }
    try {
      await ref.read(chatsApiProvider).report(_id, [message.id], reason);
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text(context.l10n.chatsReported)));
      }
    } on Object catch (error) {
      _fail(error);
    }
  }

  Future<void> _act(Chat chat, ChatMessage message, bool mine) async {
    final actions = [
      if (chat.writable) MessageAction.reply,
      if (mine && chat.writable && message.text.isNotEmpty) MessageAction.edit,
      if (!chat.archived &&
          (mine || (chat.kind != ChatKind.direct && chat.role.manages)))
        MessageAction.delete,
      if (message.text.isNotEmpty) MessageAction.copy,
      if (!mine) MessageAction.report,
    ];
    final action = await MessageActionsSheet.show(context, actions);
    if (action == null || !mounted) {
      return;
    }
    switch (action) {
      case MessageAction.reply:
        setState(() => _replyTo = message);
      case MessageAction.copy:
        await Clipboard.setData(ClipboardData(text: message.text));
        if (mounted) {
          ScaffoldMessenger.of(context)
              .showSnackBar(SnackBar(content: Text(context.l10n.chatsCopied)));
        }
      case MessageAction.edit:
        final text = await EditMessageDialog.show(context, message.text);
        if (text != null && text != message.text) {
          await _change(
            () => ref.read(chatsApiProvider).edit(message.id, text),
          );
        }
      case MessageAction.delete:
        await _change(() => ref.read(chatsApiProvider).delete(message.id));
      case MessageAction.report:
        await _report(message);
    }
  }

  Future<void> _change(Future<ChatMessage> Function() request) async {
    try {
      final changed = await request();
      await ref.read(chatsStoreProvider).saveMessages(_id, [changed]);
    } on Object catch (error) {
      _fail(error);
    }
  }

  Future<void> _accept() async {
    try {
      final changed = await ref.read(chatsApiProvider).accept(_id);
      await ref.read(chatsStoreProvider).saveChat(changed);
    } on Object catch (error) {
      _fail(error);
    }
  }

  void _fail(Object error) {
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(context.l10n, error))),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    ref.watch(chatsSyncStatusProvider);
    final chat = ref.watch(chatProvider(_id)).value;
    final entries = ref.watch(chatEntriesProvider(_id)).value ?? const [];
    final older = ref.watch(chatOlderProvider(_id)).value ?? false;
    final me = ref.watch(myIdProvider).value;
    final details = ref.watch(chatDetailsProvider(_id)).value;
    final names = {
      for (final member in details?.members ?? const <ChatMember>[])
        member.userId: member.displayName,
    };
    final bySeq = {
      for (final entry in entries)
        if (entry.message != null) entry.message!.seq: entry.message!,
    };
    ReplyQuote? quoteOf(ChatMessage? message) => message == null
        ? null
        : ReplyQuote(
            author: names[message.authorId] ?? l10n.chatsUnknown,
            text: message.deleted ? l10n.chatsDeleted : message.text,
          );
    if (chat == null) {
      return Scaffold(
        appBar: AppBar(),
        body: const Center(child: CircularProgressIndicator()),
      );
    }
    final shown = entries.reversed.toList();
    return Scaffold(
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(chat.title, overflow: TextOverflow.ellipsis),
            if (chat.kind != ChatKind.direct)
              Text(
                chat.kind == ChatKind.channel
                    ? '${l10n.chatsChannel}, ${l10n.chatsMembers(chat.memberCount)}'
                    : l10n.chatsMembers(chat.memberCount),
                style: Theme.of(context).textTheme.labelSmall,
              ),
          ],
        ),
        actions: [
          IconButton(
            key: const ValueKey('chat-info'),
            tooltip: l10n.chatsInfo,
            icon: const Icon(Icons.info_outline),
            onPressed: () => context.push(AppRoutes.chatInfo(_id)),
          ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 760),
            child: Column(
              children: [
                if (chat.request)
                  Padding(
                    padding: const EdgeInsets.all(Spacing.m),
                    child: Column(
                      children: [
                        NoticeBanner(
                          icon: Icons.mark_email_unread_outlined,
                          message: l10n.chatsRequestBanner,
                        ),
                        const SizedBox(height: Spacing.s),
                        FilledButton(
                          onPressed: _accept,
                          child: Text(l10n.chatsAccept),
                        ),
                      ],
                    ),
                  ),
                Expanded(
                  child: ListView.separated(
                    reverse: true,
                    padding: const EdgeInsets.all(Spacing.l),
                    itemCount: shown.length + (older ? 1 : 0),
                    separatorBuilder: (_, _) =>
                        const SizedBox(height: Spacing.s),
                    itemBuilder: (context, index) {
                      if (index == shown.length) {
                        return Center(
                          child: TextButton(
                            onPressed: _loadingOlder ? null : _older,
                            child: Text(l10n.chatsOlder),
                          ),
                        );
                      }
                      final entry = shown[index];
                      return _entry(chat, entry, me, names, bySeq, quoteOf);
                    },
                  ),
                ),
                if (chat.writable)
                  ChatComposer(
                    onSend: _send,
                    onAttach: _attach,
                    members: [
                      for (final member
                          in details?.members ?? const <ChatMember>[])
                        if (member.userId != me) member,
                    ],
                    quote: quoteOf(_replyTo),
                    onCancelReply: () => setState(() => _replyTo = null),
                  )
                else
                  Padding(
                    padding: const EdgeInsets.all(Spacing.l),
                    child: Text(
                      chat.archived
                          ? l10n.chatsArchived
                          : l10n.chatsReadOnlyChannel,
                      style: TextStyle(color: context.tokens.palette.muted),
                      textAlign: TextAlign.center,
                    ),
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _entry(
    Chat chat,
    ChatEntry entry,
    String? me,
    Map<String, String> names,
    Map<int, ChatMessage> bySeq,
    ReplyQuote? Function(ChatMessage?) quoteOf,
  ) {
    final outgoing = entry.outgoing;
    if (outgoing != null) {
      final store = ref.read(chatsStoreProvider);
      return OutgoingBubble(
        key: ValueKey(entry.key),
        outgoing: outgoing,
        onRetry: () => store.retry(outgoing.clientId),
        onDiscard: () => store.discard(outgoing.clientId),
      );
    }
    final message = entry.message!;
    final mine = me != null && message.authorId == me;
    final reply = message.replyTo;
    return ChatBubble(
      key: ValueKey(entry.key),
      message: message,
      mine: mine,
      author: chat.kind == ChatKind.direct
          ? null
          : names[message.authorId] ?? context.l10n.chatsUnknown,
      quote: reply == null ? null : quoteOf(bySeq[reply]),
      onLongPress: () => _act(chat, message, mine),
      onFile: _save,
    );
  }
}
