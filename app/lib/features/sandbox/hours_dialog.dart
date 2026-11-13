import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'hours_request.dart';

class HoursDialog extends StatefulWidget {
  const HoursDialog({
    required this.title,
    required this.action,
    this.hours,
    this.withReason = false,
    super.key,
  });

  final String title;
  final String action;
  final double? hours;
  final bool withReason;

  static Future<HoursRequest?> show(
    BuildContext context, {
    required String title,
    required String action,
    double? hours,
    bool withReason = false,
  }) => showDialog<HoursRequest>(
    context: context,
    builder: (context) => HoursDialog(
      title: title,
      action: action,
      hours: hours,
      withReason: withReason,
    ),
  );

  @override
  State<HoursDialog> createState() => _HoursDialogState();
}

class _HoursDialogState extends State<HoursDialog> {
  late final _hours = TextEditingController(
    text: widget.hours == null ? '' : _plain(widget.hours!),
  );
  final _reason = TextEditingController();

  static String _plain(double value) => value == value.roundToDouble()
      ? value.toStringAsFixed(0)
      : value.toString();

  @override
  void dispose() {
    _hours.dispose();
    _reason.dispose();
    super.dispose();
  }

  void _done() {
    final hours = double.tryParse(_hours.text.trim().replaceAll(',', '.'));
    if (hours == null || hours <= 0) {
      return;
    }
    Navigator.of(context)
        .pop(HoursRequest(hours: hours, reason: _reason.text.trim()));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(widget.title),
      content: SizedBox(
        width: 360,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              key: const ValueKey('sandbox-hours'),
              controller: _hours,
              autofocus: true,
              keyboardType: const TextInputType.numberWithOptions(
                decimal: true,
              ),
              decoration: InputDecoration(labelText: l10n.sandboxHoursField),
            ),
            if (widget.withReason)
              TextField(
                key: const ValueKey('sandbox-reason'),
                controller: _reason,
                maxLines: 3,
                decoration: InputDecoration(labelText: l10n.sandboxReasonField),
              ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(
          key: const ValueKey('sandbox-hours-done'),
          onPressed: _done,
          child: Text(widget.action),
        ),
      ],
    );
  }
}
