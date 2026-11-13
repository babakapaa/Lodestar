import 'package:material_ui/material_ui.dart';

import '../../design/components/choice_group.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/assignment.dart';
import '../study/grade_input.dart';
import '../study/grading.dart';
import '../study/study_labels_context.dart';
import 'defense_labels.dart';
import 'defense_status.dart';
import 'finish_choice.dart';

class FinishDefenseDialog extends StatefulWidget {
  const FinishDefenseDialog({required this.assignment, super.key});

  final Assignment assignment;

  static Future<FinishChoice?> show(
    BuildContext context,
    Assignment assignment,
  ) => showDialog<FinishChoice>(
    context: context,
    builder: (context) => FinishDefenseDialog(assignment: assignment),
  );

  @override
  State<FinishDefenseDialog> createState() => _FinishDefenseDialogState();
}

class _FinishDefenseDialogState extends State<FinishDefenseDialog> {
  DefenseStatus _status = DefenseStatus.held;
  bool _grade = true;
  bool _passed = true;
  final _points = TextEditingController();
  final _comment = TextEditingController();
  bool _invalid = false;

  @override
  void dispose() {
    _points.dispose();
    _comment.dispose();
    super.dispose();
  }

  GradeInput? _input() {
    final a = widget.assignment;
    final comment = _comment.text.trim();
    switch (a.grading) {
      case Grading.points:
        final points = double.tryParse(
          _points.text.trim().replaceAll(',', '.').replaceAll(' ', ''),
        );
        if (points == null || points < 0 || points > (a.maxPoints ?? 0)) {
          return null;
        }
        return GradeInput(points: points, passed: null, comment: comment);
      case Grading.pass:
        return GradeInput(points: null, passed: _passed, comment: comment);
      case Grading.none:
        return comment.isEmpty
            ? null
            : GradeInput(points: null, passed: null, comment: comment);
    }
  }

  void _done() {
    GradeInput? grade;
    if (_grade) {
      grade = _input();
      if (grade == null) {
        setState(() => _invalid = true);
        return;
      }
    }
    Navigator.of(context).pop(FinishChoice(status: _status, grade: grade));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = DefenseLabels(l10n);
    final a = widget.assignment;
    return AlertDialog(
      title: Text(l10n.defenseFinish),
      content: SizedBox(
        width: 440,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              ChoiceGroup<DefenseStatus>(
                values: const [DefenseStatus.held, DefenseStatus.missed],
                selected: _status,
                labelOf: labels.status,
                onSelected: (status) => setState(() => _status = status),
              ),
              SwitchListTile(
                key: const ValueKey('finish-grade'),
                contentPadding: EdgeInsets.zero,
                title: Text(l10n.defenseGradeNow),
                value: _grade,
                onChanged: (value) => setState(() => _grade = value),
              ),
              if (_grade) ...[
                switch (a.grading) {
                  Grading.points => TextField(
                    key: const ValueKey('finish-points'),
                    controller: _points,
                    keyboardType: const TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    decoration: InputDecoration(
                      labelText: l10n.gradePointsLabel(
                        context.studyLabels.points(a.maxPoints ?? 0),
                      ),
                    ),
                  ),
                  Grading.pass => ChoiceGroup<bool>(
                    values: const [true, false],
                    selected: _passed,
                    labelOf: (passed) =>
                        passed ? l10n.gradePassed : l10n.gradeFailed,
                    onSelected: (value) => setState(() => _passed = value),
                  ),
                  Grading.none => const SizedBox.shrink(),
                },
                const SizedBox(height: Spacing.m),
                TextField(
                  key: const ValueKey('finish-comment'),
                  controller: _comment,
                  minLines: 2,
                  maxLines: 5,
                  decoration: InputDecoration(labelText: l10n.gradeComment),
                ),
              ],
              if (_invalid)
                Text(
                  l10n.errorInvalidGrade,
                  style: TextStyle(color: Theme.of(context).colorScheme.error),
                ),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(onPressed: _done, child: Text(l10n.defenseFinishSave)),
      ],
    );
  }
}
