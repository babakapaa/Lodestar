import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'attachment_chips.dart';
import 'chat_draft.dart';
import 'chat_file.dart';
import 'chat_member.dart';
import 'mention_query.dart';
import 'mention_suggestions.dart';
import 'reply_quote.dart';

class ChatComposer extends StatefulWidget {
  const ChatComposer({
    required this.onSend,
    this.members = const [],
    this.onAttach,
    this.quote,
    this.onCancelReply,
    super.key,
  });

  static const maxFiles = 10;

  final ValueChanged<ChatDraft> onSend;
  final List<ChatMember> members;
  final Future<List<ChatFile>> Function()? onAttach;
  final ReplyQuote? quote;
  final VoidCallback? onCancelReply;

  @override
  State<ChatComposer> createState() => _ChatComposerState();
}

class _ChatComposerState extends State<ChatComposer> {
  final _text = TextEditingController();
  final _mentions = <String, String>{};
  final _files = <ChatFile>[];
  MentionQuery? _query;
  bool _attaching = false;

  @override
  void initState() {
    super.initState();
    _text.addListener(_track);
  }

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  void _track() => setState(() => _query = MentionQuery.of(_text.value));

  List<ChatMember> get _suggestions {
    final query = _query;
    if (query == null) {
      return const [];
    }
    final needle = query.text.toLowerCase();
    return [
      for (final member in widget.members)
        if (member.displayName.toLowerCase().contains(needle)) member,
    ].take(6).toList();
  }

  void _mention(ChatMember member) {
    final query = _query;
    if (query == null) {
      return;
    }
    _mentions[member.userId] = member.displayName;
    _text.value = query.insert(_text.value, member.displayName);
  }

  Future<void> _attach() async {
    final attach = widget.onAttach;
    if (attach == null) {
      return;
    }
    setState(() => _attaching = true);
    try {
      final files = await attach();
      if (mounted) {
        setState(() {
          _files.addAll(files.take(ChatComposer.maxFiles - _files.length));
        });
      }
    } finally {
      if (mounted) {
        setState(() => _attaching = false);
      }
    }
  }

  void _send() {
    final text = _text.text.trim();
    if ((text.isEmpty && _files.isEmpty) || _attaching) {
      return;
    }
    final draft = ChatDraft(
      text: text,
      mentions: [
        for (final MapEntry(:key, :value) in _mentions.entries)
          if (text.contains('@$value')) key,
      ],
      files: List.of(_files),
    );
    _text.clear();
    setState(() {
      _mentions.clear();
      _files.clear();
    });
    widget.onSend(draft);
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final quote = widget.quote;
    final suggestions = _suggestions;
    return Padding(
      padding: const EdgeInsets.all(Spacing.m),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          if (suggestions.isNotEmpty) ...[
            MentionSuggestions(members: suggestions, onPick: _mention),
            const SizedBox(height: Spacing.s),
          ],
          if (quote != null)
            Row(
              children: [
                Expanded(child: quote),
                IconButton(
                  tooltip: l10n.actionCancel,
                  icon: const Icon(Icons.close),
                  onPressed: widget.onCancelReply,
                ),
              ],
            ),
          if (_files.isNotEmpty) ...[
            AttachmentChips(
              files: _files,
              onRemove: (file) => setState(() => _files.remove(file)),
            ),
            const SizedBox(height: Spacing.s),
          ],
          Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              if (widget.onAttach != null)
                IconButton(
                  key: const ValueKey('chat-attach'),
                  tooltip: l10n.chatsAttach,
                  onPressed:
                      _attaching || _files.length >= ChatComposer.maxFiles
                      ? null
                      : _attach,
                  icon: _attaching
                      ? const SizedBox.square(
                          dimension: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Icon(Icons.attach_file),
                ),
              Expanded(
                child: TextField(
                  key: const ValueKey('chat-message'),
                  controller: _text,
                  minLines: 1,
                  maxLines: 6,
                  maxLength: 4000,
                  textInputAction: TextInputAction.send,
                  onSubmitted: (_) => _send(),
                  decoration: InputDecoration(
                    hintText: l10n.chatsMessageHint,
                    counterText: '',
                  ),
                ),
              ),
              const SizedBox(width: Spacing.s),
              IconButton.filled(
                tooltip: l10n.chatsSend,
                onPressed: _send,
                icon: const Icon(Icons.send),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
