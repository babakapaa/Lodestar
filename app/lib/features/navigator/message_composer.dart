import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class MessageComposer extends StatefulWidget {
  const MessageComposer({
    required this.enabled,
    required this.onSend,
    this.initial = '',
    super.key,
  });

  final bool enabled;
  final String initial;
  final ValueChanged<String> onSend;

  @override
  State<MessageComposer> createState() => _MessageComposerState();
}

class _MessageComposerState extends State<MessageComposer> {
  late final _text = TextEditingController(text: widget.initial);

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
    final palette = context.palette;
    return Padding(
      padding: const EdgeInsets.fromLTRB(
        Spacing.l,
        Spacing.s,
        Spacing.l,
        Spacing.l,
      ),
      child: Container(
        padding: const EdgeInsets.fromLTRB(
          Spacing.xs,
          Spacing.xs,
          Spacing.s,
          Spacing.s,
        ),
        decoration: BoxDecoration(
          color: palette.surface,
          borderRadius: BorderRadius.circular(20),
          border: Border.all(color: palette.border),
          boxShadow: [
            BoxShadow(
              color: palette.text.withValues(alpha: 0.06),
              blurRadius: 12,
              offset: const Offset(0, 4),
            ),
          ],
        ),
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
                  filled: false,
                  border: InputBorder.none,
                  enabledBorder: InputBorder.none,
                  focusedBorder: InputBorder.none,
                ),
              ),
            ),
            const SizedBox(width: Spacing.s),
            IconButton.filled(
              tooltip: l10n.navigatorSend,
              onPressed: widget.enabled ? _send : null,
              icon: const Icon(Icons.arrow_upward),
            ),
          ],
        ),
      ),
    );
  }
}
