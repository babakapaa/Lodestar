import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class ReportDialog extends StatefulWidget {
  const ReportDialog({super.key});

  static Future<String?> show(BuildContext context) => showDialog<String>(
    context: context,
    builder: (context) => const ReportDialog(),
  );

  @override
  State<ReportDialog> createState() => _ReportDialogState();
}

class _ReportDialogState extends State<ReportDialog> {
  final _reason = TextEditingController();

  @override
  void dispose() {
    _reason.dispose();
    super.dispose();
  }

  void _done() {
    final reason = _reason.text.trim();
    if (reason.isNotEmpty) {
      Navigator.of(context).pop(reason);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.chatsReportTitle),
      content: SizedBox(
        width: 480,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(l10n.chatsReportHint),
            TextField(
              key: const ValueKey('chat-report-reason'),
              controller: _reason,
              autofocus: true,
              minLines: 2,
              maxLines: 6,
              maxLength: 1000,
              decoration: InputDecoration(labelText: l10n.chatsReportReason),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(onPressed: _done, child: Text(l10n.chatsReportSend)),
      ],
    );
  }
}
