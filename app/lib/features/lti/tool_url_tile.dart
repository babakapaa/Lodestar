import 'package:flutter/services.dart';
import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class ToolUrlTile extends StatelessWidget {
  const ToolUrlTile({required this.label, required this.url, super.key});

  final String label;
  final String url;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return ListTile(
      contentPadding: EdgeInsets.zero,
      title: Text(label),
      subtitle: SelectableText(url),
      trailing: IconButton(
        tooltip: l10n.lmsCopy,
        icon: const Icon(Icons.copy_outlined),
        onPressed: () async {
          final messenger = ScaffoldMessenger.of(context);
          await Clipboard.setData(ClipboardData(text: url));
          messenger.showSnackBar(SnackBar(content: Text(l10n.lmsCopied)));
        },
      ),
    );
  }
}
