import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class ConsentJoinDialog extends StatelessWidget {
  const ConsentJoinDialog({super.key});

  static Future<bool> show(BuildContext context) async =>
      await showDialog<bool>(
        context: context,
        builder: (context) => const ConsentJoinDialog(),
      ) ??
      false;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.callRecording),
      content: Text(l10n.callConsentJoin),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(false),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(
          key: const ValueKey('consent-join'),
          onPressed: () => Navigator.of(context).pop(true),
          child: Text(l10n.callConsentAgree),
        ),
      ],
    );
  }
}
