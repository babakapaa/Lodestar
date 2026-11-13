import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'estimate_option.dart';
import 'run_estimate.dart';
import 'sandbox_labels_context.dart';

class ClassOptions extends StatelessWidget {
  const ClassOptions({
    required this.estimate,
    required this.selected,
    required this.onSelected,
    super.key,
  });

  final RunEstimate estimate;
  final String? selected;
  final ValueChanged<String> onSelected;

  String _details(BuildContext context, EstimateOption o) {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    return [
      if (estimate.starter)
        l10n.sandboxFree
      else
        l10n.sandboxUpTo(labels.hours(o.cost)),
      if (estimate.night != null && !estimate.starter)
        l10n.sandboxAtNight(labels.hours(o.nightCost)),
      if (o.startsNow)
        l10n.sandboxStartsNow
      else
        l10n.sandboxQueueLength(o.queued),
      if (!o.enough) l10n.sandboxNotEnough,
    ].join(', ');
  }

  @override
  Widget build(BuildContext context) {
    return RadioGroup<String>(
      groupValue: selected,
      onChanged: (value) {
        if (value != null) {
          onSelected(value);
        }
      },
      child: Column(
        children: [
          for (final o in estimate.options)
            RadioListTile<String>(
              key: ValueKey('class-${o.classId}'),
              contentPadding: EdgeInsets.zero,
              value: o.classId,
              enabled: o.enough,
              title: Text(
                context.l10n.sandboxClassName(o.className, o.memoryGb),
              ),
              subtitle: Text(_details(context, o)),
            ),
        ],
      ),
    );
  }
}
