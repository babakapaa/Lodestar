import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';

class MessageComposer extends StatefulWidget {
  const MessageComposer({
    required this.enabled,
    required this.onSend,
    super.key,
  });

  final bool enabled;
  final ValueChanged<String> onSend;

  @override
  State<MessageComposer> createState() => _MessageComposerState();
}

class _MessageComposerState extends State<MessageComposer> {
  final _text = TextEditingController();

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  void _send() {
    final text = _text.text.trim();
    if (text.isEmpty || !widget.enabled) {
      return;
    }
    _text.clear();
    widget.onSend(text);
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return Padding(
      padding: const EdgeInsets.all(Spacing.m),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Expanded(
            child: TextField(
              key: const ValueKey('navigator-question'),
              controller: _text,
              minLines: 1,
              maxLines: 6,
              maxLength: 8000,
              textInputAction: TextInputAction.send,
              onSubmitted: (_) => _send(),
              decoration: InputDecoration(
                hintText: l10n.navigatorAskHint,
                counterText: '',
              ),
            ),
          ),
          const SizedBox(width: Spacing.s),
          IconButton.filled(
            tooltip: l10n.navigatorSend,
            onPressed: widget.enabled ? _send : null,
            icon: const Icon(Icons.send),
          ),
        ],
      ),
    );
  }
}
