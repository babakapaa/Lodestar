import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'grade_parts.dart';
import 'grade_weights.dart';
import 'points_text.dart';
import 'study_labels_context.dart';

class GradePartsFields extends StatefulWidget {
  const GradePartsFields({
    required this.weights,
    required this.maxPoints,
    required this.initial,
    required this.onChanged,
    super.key,
  });

  final GradeWeights weights;
  final double maxPoints;
  final GradeParts? initial;
  final ValueChanged<GradeParts> onChanged;

  @override
  State<GradePartsFields> createState() => _GradePartsFieldsState();
}

class _GradePartsFieldsState extends State<GradePartsFields> {
  late final _defense = _controller(widget.initial?.defense);
  late final _map = _controller(widget.initial?.map);
  late final _work = _controller(widget.initial?.work);

  static TextEditingController _controller(double? value) =>
      TextEditingController(text: value == null ? '' : PointsText.of(value));

  @override
  void dispose() {
    _defense.dispose();
    _map.dispose();
    _work.dispose();
    super.dispose();
  }

  static double? _value(TextEditingController controller, int weight) =>
      weight == 0
      ? null
      : double.tryParse(controller.text.trim().replaceAll(',', '.'));

  GradeParts get _parts => GradeParts(
    defense: _value(_defense, widget.weights.defense),
    map: _value(_map, widget.weights.map),
    work: _value(_work, widget.weights.work),
  );

  void _emit() {
    setState(() {});
    widget.onChanged(_parts);
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final weights = widget.weights;
    final points = weights.points(_parts, widget.maxPoints);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        for (final (controller, weight, label, key) in [
          (_defense, weights.defense, l10n.gradePartDefense, 'part-defense'),
          (_map, weights.map, l10n.gradePartMap, 'part-map'),
          (_work, weights.work, l10n.gradePartWork, 'part-work'),
        ])
          if (weight > 0)
            Padding(
              padding: const EdgeInsets.only(bottom: Spacing.s),
              child: TextField(
                key: ValueKey(key),
                controller: controller,
                keyboardType: const TextInputType.numberWithOptions(
                  decimal: true,
                ),
                decoration: InputDecoration(
                  labelText: l10n.gradePartLabel(label, weight),
                  suffixText: '%',
                ),
                onChanged: (_) => _emit(),
              ),
            ),
        Text(
          points == null
              ? l10n.gradePartsHint
              : l10n.gradePartsTotal(
                  labels.points(points),
                  labels.points(widget.maxPoints),
                ),
          style: Theme.of(context).textTheme.bodyMedium!
              .copyWith(color: context.palette.muted),
        ),
      ],
    );
  }
}
