import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class StatsPeriod extends StatelessWidget {
  const StatsPeriod({required this.days, required this.onChanged, super.key});

  final int days;
  final ValueChanged<int> onChanged;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return SegmentedButton<int>(
      showSelectedIcon: false,
      segments: [
        ButtonSegment(value: 7, label: Text(l10n.statsWeek)),
        ButtonSegment(value: 30, label: Text(l10n.statsMonth)),
        ButtonSegment(value: 90, label: Text(l10n.statsQuarter)),
        ButtonSegment(value: 365, label: Text(l10n.statsYear)),
      ],
      selected: {days},
      onSelectionChanged: (selected) => onChanged(selected.first),
    );
  }
}
