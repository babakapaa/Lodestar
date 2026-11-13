import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class TeamTextDialog extends StatefulWidget {
  const TeamTextDialog({
    required this.title,
    required this.label,
    required this.action,
    this.initial = '',
    this.hint,
    this.required = true,
    this.lines = 1,
    super.key,
  });

  final String title;
  final String label;
  final String action;
  final String initial;
  final String? hint;
  final bool required;
  final int lines;

  static Future<String?> show(
    BuildContext context, {
    required String title,
    required String label,
    required String action,
    String initial = '',
    String? hint,
    bool required = true,
    int lines = 1,
  }) => showDialog<String>(
    context: context,
    builder: (context) => TeamTextDialog(
      title: title,
      label: label,
      action: action,
      initial: initial,
      hint: hint,
      required: required,
      lines: lines,
    ),
  );

  @override
  State<TeamTextDialog> createState() => _TeamTextDialogState();
}

class _TeamTextDialogState extends State<TeamTextDialog> {
  late final _text = TextEditingController(text: widget.initial);

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  void _done() {
    final value = _text.text.trim();
    if (widget.required && value.isEmpty) {
      return;
    }
    Navigator.of(context).pop(value);
  }

  @override
  Widget build(BuildContext context) {
    final hint = widget.hint;
    return AlertDialog(
      title: Text(widget.title),
      content: SizedBox(
        width: 380,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (hint != null) Text(hint),
            TextField(
              key: const ValueKey('team-text'),
              controller: _text,
              autofocus: true,
              minLines: widget.lines,
              maxLines: widget.lines == 1 ? 1 : widget.lines + 4,
              decoration: InputDecoration(labelText: widget.label),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(context.l10n.actionCancel),
        ),
        FilledButton(onPressed: _done, child: Text(widget.action)),
      ],
    );
  }
}
