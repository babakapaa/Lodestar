import 'package:flutter/services.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';

class NodeTokenDialog extends StatelessWidget {
  const NodeTokenDialog({required this.token, super.key});

  final String token;

  static Future<void> show(BuildContext context, String token) =>
      showDialog<void>(
        context: context,
        builder: (context) => NodeTokenDialog(token: token),
      );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.sandboxNodeToken),
      content: SizedBox(
        width: 420,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(l10n.sandboxNodeTokenHint),
            const SizedBox(height: Spacing.m),
            SelectableText(
              token,
              key: const ValueKey('node-token'),
              style: const TextStyle(fontFamily: 'monospace'),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Clipboard.setData(ClipboardData(text: token)),
          child: Text(l10n.sandboxCopy),
        ),
        FilledButton(
          key: const ValueKey('node-token-done'),
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.sandboxDone),
        ),
      ],
    );
  }
}
