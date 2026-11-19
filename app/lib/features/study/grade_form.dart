import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/choice_group.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'assignment.dart';
import 'grade.dart';
import 'grade_input.dart';
import 'grade_parts.dart';
import 'grade_parts_fields.dart';
import 'grading.dart';
import 'points_text.dart';
import 'study_labels_context.dart';
import 'study_providers.dart';

class GradeForm extends ConsumerStatefulWidget {
  const GradeForm({
    required this.assignment,
    required this.studentId,
    required this.grade,
    super.key,
  });

  final Assignment assignment;
  final String studentId;
  final Grade? grade;

  @override
  ConsumerState<GradeForm> createState() => _GradeFormState();
}

class _GradeFormState extends ConsumerState<GradeForm> {
  late final _points = TextEditingController(
    text: switch (widget.grade?.points) {
      final points? => PointsText.of(points),
      null => '',
    },
  );
  late final _comment = TextEditingController(
    text: widget.grade?.comment ?? '',
  );
  late bool _passed = widget.grade?.passed ?? true;
  late bool _byParts =
      widget.assignment.gradeWeights != null &&
      (widget.grade == null || widget.grade?.parts != null);
  late GradeParts _parts = widget.grade?.parts ?? const GradeParts();
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    _points.dispose();
    _comment.dispose();
    super.dispose();
  }

  GradeInput? _input() {
    final assignment = widget.assignment;
    final comment = _comment.text.trim();
    switch (assignment.grading) {
      case Grading.points:
        final weights = assignment.gradeWeights;
        if (_byParts && weights != null) {
          if (weights.points(_parts, assignment.maxPoints ?? 0) == null) {
            return null;
          }
          return GradeInput(
            points: null,
            passed: null,
            comment: comment,
            parts: _parts,
          );
        }
        final points = double.tryParse(
          _points.text.trim().replaceAll(',', '.').replaceAll(' ', ''),
        );
        final max = assignment.maxPoints ?? 0;
        if (points == null || points < 0 || points > max) {
          return null;
        }
        return GradeInput(points: points, passed: null, comment: comment);
      case Grading.pass:
        return GradeInput(points: null, passed: _passed, comment: comment);
      case Grading.none:
        if (comment.isEmpty) {
          return null;
        }
        return GradeInput(points: null, passed: null, comment: comment);
    }
  }

  Future<void> _save() async {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    final input = _input();
    if (input == null) {
      setState(() => _error = l10n.errorInvalidGrade);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref
          .read(studyRepositoryProvider)
          .grade(widget.assignment.id, widget.studentId, input);
      ref
        ..invalidate(assignmentWorkProvider(widget.assignment.id))
        ..invalidate(assignmentProvider(widget.assignment.id))
        ..invalidate(assignmentsProvider(widget.assignment.courseId));
      router.pop();
    } on ApiException catch (error) {
      if (mounted) {
        setState(() => _error = ApiErrorText.of(l10n, error));
      }
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final assignment = widget.assignment;
    final error = _error;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (assignment.grading == Grading.points &&
            assignment.gradeWeights != null) ...[
          SwitchListTile(
            key: const ValueKey('grade-by-parts'),
            contentPadding: EdgeInsets.zero,
            title: Text(l10n.gradeByParts),
            value: _byParts,
            onChanged: (value) => setState(() => _byParts = value),
          ),
          const SizedBox(height: Spacing.s),
        ],
        switch (assignment.grading) {
          Grading.points when _byParts && assignment.gradeWeights != null =>
            GradePartsFields(
              weights: assignment.gradeWeights!,
              maxPoints: assignment.maxPoints ?? 0,
              initial: widget.grade?.parts,
              onChanged: (value) => _parts = value,
            ),
          Grading.points => TextField(
            controller: _points,
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
            decoration: InputDecoration(
              labelText: l10n.gradePointsLabel(
                labels.points(assignment.maxPoints ?? 0),
              ),
            ),
          ),
          Grading.pass => ChoiceGroup<bool>(
            values: const [true, false],
            selected: _passed,
            labelOf: (passed) => passed ? l10n.gradePassed : l10n.gradeFailed,
            onSelected: (value) => setState(() => _passed = value),
          ),
          Grading.none => const SizedBox.shrink(),
        },
        const SizedBox(height: Spacing.l),
        TextField(
          controller: _comment,
          minLines: 3,
          maxLines: 8,
          decoration: InputDecoration(
            labelText: l10n.gradeComment,
            alignLabelWithHint: true,
          ),
        ),
        if (error != null) FormError(message: error),
        const SizedBox(height: Spacing.l),
        Align(
          alignment: Alignment.centerLeft,
          child: LodestarButton(
            label: l10n.gradeSave,
            onPressed: _busy ? null : _save,
          ),
        ),
      ],
    );
  }
}
