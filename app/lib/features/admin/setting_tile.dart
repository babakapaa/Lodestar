import 'package:material_ui/material_ui.dart';

import '../../core/settings/preference_value.dart';
import '../../l10n/l10n_context.dart';
import 'setting_definition.dart';
import 'setting_labels.dart';
import 'setting_value_dialog.dart';

class SettingTile extends StatelessWidget {
  const SettingTile({
    required this.definition,
    required this.value,
    required this.onChanged,
    super.key,
  });

  final SettingDefinition definition;
  final PreferenceValue? value;
  final ValueChanged<Object> onChanged;

  Object? get _current => value?.value ?? definition.defaultValue;

  Future<void> _editText(BuildContext context, String label) async {
    final saved = await showDialog<String>(
      context: context,
      builder: (context) => SettingValueDialog(
        definition: definition,
        title: label,
        initial: '${_current ?? ''}',
      ),
    );
    if (saved == null) {
      return;
    }
    if (definition.type != SettingDefinition.integerType) {
      onChanged(saved.trim());
      return;
    }
    final number = int.tryParse(saved);
    if (number != null) {
      onChanged(number);
    }
  }

  @override
  Widget build(BuildContext context) {
    final labels = SettingLabels(context.l10n);
    final label = labels.name(definition.key);
    final current = _current;
    return switch (definition.type) {
      SettingDefinition.boolType => SwitchListTile(
        contentPadding: EdgeInsets.zero,
        title: Text(label),
        value: current == true,
        onChanged: onChanged,
      ),
      SettingDefinition.choiceType => ListTile(
        contentPadding: EdgeInsets.zero,
        title: Text(label),
        trailing: DropdownButton<String>(
          value: definition.options.contains(current)
              ? current! as String
              : null,
          items: [
            for (final option in definition.options)
              DropdownMenuItem(
                value: option,
                child: Text(labels.option(definition.key, option)),
              ),
          ],
          onChanged: (option) {
            if (option != null) {
              onChanged(option);
            }
          },
        ),
      ),
      _ => ListTile(
        contentPadding: EdgeInsets.zero,
        title: Text(label),
        trailing: Text('${current ?? ''}'),
        onTap: () => _editText(context, label),
      ),
    };
  }
}
