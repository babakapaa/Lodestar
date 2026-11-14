import 'package:material_ui/material_ui.dart';

import '../../design/tokens/accent_choice.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/theme_family.dart';
import 'accent_swatch.dart';

class AccentPicker extends StatelessWidget {
  const AccentPicker({
    required this.family,
    required this.selected,
    required this.labelOf,
    required this.onSelected,
    super.key,
  });

  final ThemeFamily family;
  final AccentChoice selected;
  final String Function(AccentChoice) labelOf;
  final ValueChanged<AccentChoice> onSelected;

  @override
  Widget build(BuildContext context) {
    final brightness = Theme.of(context).brightness;
    final base = family.spec.palette(brightness);
    return Wrap(
      spacing: Spacing.m,
      runSpacing: Spacing.m,
      children: [
        for (final choice in AccentChoice.values)
          AccentSwatch(
            color: choice.apply(base, brightness).accent,
            label: labelOf(choice),
            auto: choice == AccentChoice.auto,
            selected: choice == selected,
            onTap: () => onSelected(choice),
          ),
      ],
    );
  }
}
