import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'reply_quote.dart';

class ChatComposer extends StatefulWidget {
  const ChatComposer({
    required this.onSend,
    this.quote,
    this.onCancelReply,
    super.key,
  });

  final ValueChanged<String> onSend;
  final ReplyQuote? quote;
  final VoidCallback? onCancelReply;

  @override
  State<ChatComposer> createState() => _ChatComposerState();
}

class _ChatComposerState extends State<ChatComposer> {
  final _text = TextEditingController();

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  void _send() {
    final text = _text.text.trim();
    if (text.isEmpty) {
      return;
    }
    _text.clear();
    widget.onSend(text);
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final quote = widget.quote;
    return Padding(
      padding: const EdgeInsets.all(Spacing.m),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
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
          Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
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
