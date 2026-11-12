import 'package:flutter/services.dart';
import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'setting_definition.dart';

class SettingValueDialog extends StatefulWidget {
  const SettingValueDialog({
    required this.definition,
    required this.title,
    required this.initial,
    super.key,
  });

  final SettingDefinition definition;
  final String title;
  final String initial;

  @override
  State<SettingValueDialog> createState() => _SettingValueDialogState();
}

class _SettingValueDialogState extends State<SettingValueDialog> {
  late final _controller = TextEditingController(text: widget.initial);

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final definition = widget.definition;
    final integer = definition.type == SettingDefinition.integerType;
    final min = definition.min;
    final max = definition.max;
    return AlertDialog(
      title: Text(widget.title),
      content: TextField(
        controller: _controller,
        autofocus: true,
        keyboardType: integer ? TextInputType.number : TextInputType.text,
        inputFormatters: [if (integer) FilteringTextInputFormatter.digitsOnly],
        maxLength: definition.maxLength,
        decoration: InputDecoration(
          helperText: integer && min != null && max != null
              ? l10n.settingRange(min, max)
              : null,
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        TextButton(
          onPressed: () => Navigator.of(context).pop(_controller.text),
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
