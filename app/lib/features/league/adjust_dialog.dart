import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class AdjustDialog extends StatefulWidget {
  const AdjustDialog({super.key});

  static Future<(int, String)?> show(BuildContext context) =>
      showDialog<(int, String)>(
        context: context,
        builder: (context) => const AdjustDialog(),
      );

  @override
  State<AdjustDialog> createState() => _AdjustDialogState();
}

class _AdjustDialogState extends State<AdjustDialog> {
  final _points = TextEditingController();
  final _reason = TextEditingController();

  @override
  void dispose() {
    _points.dispose();
    _reason.dispose();
    super.dispose();
  }

  void _done() {
    final points = int.tryParse(_points.text.trim());
    final reason = _reason.text.trim();
    if (points == null || points == 0 || reason.isEmpty) {
      return;
    }
    Navigator.of(context).pop((points, reason));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.leagueAdjust),
      content: SizedBox(
        width: 380,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              key: const ValueKey('adjust-points'),
              controller: _points,
              keyboardType: const TextInputType.numberWithOptions(signed: true),
              decoration: InputDecoration(labelText: l10n.leagueAdjustPoints),
            ),
            TextField(
              key: const ValueKey('adjust-reason'),
              controller: _reason,
              decoration: InputDecoration(labelText: l10n.leagueAdjustReason),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(onPressed: _done, child: Text(l10n.actionSave)),
      ],
    );
  }
}
