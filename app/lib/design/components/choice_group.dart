import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import 'lodestar_chip.dart';

class ChoiceGroup<T> extends StatelessWidget {
  const ChoiceGroup({
    required this.values,
    required this.selected,
    required this.labelOf,
    required this.onSelected,
    this.leadingOf,
    super.key,
  });

  final List<T> values;
  final T selected;
  final String Function(T value) labelOf;
  final ValueChanged<T> onSelected;
  final Widget Function(T value)? leadingOf;

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: Spacing.s,
      runSpacing: Spacing.s,
      children: [
        for (final value in values)
          LodestarChip(
            label: labelOf(value),
            selected: value == selected,
            leading: leadingOf?.call(value),
            onSelected: () => onSelected(value),
          ),
      ],
    );
  }
}
