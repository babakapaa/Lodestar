import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'grade_weights.dart';

class GradeWeightsEditor extends StatefulWidget {
  const GradeWeightsEditor({
    required this.initial,
    required this.onChanged,
    super.key,
  });

  final GradeWeights? initial;
  final ValueChanged<GradeWeights?> onChanged;

  @override
  State<GradeWeightsEditor> createState() => _GradeWeightsEditorState();
}

class _GradeWeightsEditorState extends State<GradeWeightsEditor> {
  late bool _enabled = widget.initial != null;
  late final _defense = TextEditingController(
    text: '${(widget.initial ?? GradeWeights.suggested).defense}',
  );
  late final _map = TextEditingController(
    text: '${(widget.initial ?? GradeWeights.suggested).map}',
  );
  late final _work = TextEditingController(
    text: '${(widget.initial ?? GradeWeights.suggested).work}',
  );

  @override
  void dispose() {
    _defense.dispose();
    _map.dispose();
    _work.dispose();
    super.dispose();
  }

  GradeWeights get _weights => GradeWeights(
    defense: int.tryParse(_defense.text.trim()) ?? -1,
    map: int.tryParse(_map.text.trim()) ?? -1,
    work: int.tryParse(_work.text.trim()) ?? -1,
  );

  void _emit() {
    setState(() {});
    widget.onChanged(_enabled ? _weights : null);
  }

  Widget _field(TextEditingController controller, String label, Key key) =>
      Expanded(
        child: TextField(
          key: key,
          controller: controller,
          keyboardType: TextInputType.number,
          decoration: InputDecoration(labelText: label, suffixText: '%'),
          onChanged: (_) => _emit(),
        ),
      );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        SwitchListTile(
          key: const ValueKey('assignment-weights'),
          contentPadding: EdgeInsets.zero,
          title: Text(l10n.assignmentWeights),
          subtitle: Text(l10n.assignmentWeightsHint),
          value: _enabled,
          onChanged: (value) {
            _enabled = value;
            _emit();
          },
        ),
        if (_enabled) ...[
          Row(
            children: [
              _field(
                _defense,
                l10n.gradePartDefense,
                const ValueKey('weight-defense'),
              ),
              const SizedBox(width: Spacing.m),
              _field(_map, l10n.gradePartMap, const ValueKey('weight-map')),
              const SizedBox(width: Spacing.m),
              _field(_work, l10n.gradePartWork, const ValueKey('weight-work')),
            ],
          ),
          if (!_weights.valid) ...[
            const SizedBox(height: Spacing.xs),
            Text(
              l10n.assignmentWeightsSum,
              style: TextStyle(color: Theme.of(context).colorScheme.error),
            ),
          ],
        ],
      ],
    );
  }
}
