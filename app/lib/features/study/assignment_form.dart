import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/choice_group.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'assignment_details.dart';
import 'assignment_input.dart';
import 'grade_weights.dart';
import 'grade_weights_editor.dart';
import 'grading.dart';
import 'map_rules.dart';
import 'defense_format.dart';
import 'navigator_mode.dart';
import 'points_text.dart';
import 'study_labels_context.dart';
import 'study_providers.dart';
import 'team_rules.dart';
import 'team_rules_editor.dart';

class AssignmentForm extends ConsumerStatefulWidget {
  const AssignmentForm({required this.courseId, this.details, super.key});

  final String courseId;
  final AssignmentDetails? details;

  @override
  ConsumerState<AssignmentForm> createState() => _AssignmentFormState();
}

class _AssignmentFormState extends ConsumerState<AssignmentForm> {
  late final _assignment = widget.details?.assignment;
  late final _title = TextEditingController(text: _assignment?.title ?? '');
  late final _description = TextEditingController(
    text: _assignment?.description ?? '',
  );
  late final _maxPoints = TextEditingController(
    text: PointsText.of(_assignment?.maxPoints ?? 100),
  );
  late DateTime? _dueAt = _assignment?.dueAt?.toLocal();
  late Grading _grading = _assignment?.grading ?? Grading.points;
  late bool _published = _assignment?.published ?? false;
  late GradeWeights? _weights = _assignment == null
      ? GradeWeights.suggested
      : _assignment.gradeWeights;
  late bool _mapRequired = _assignment?.searchMap != null;
  late NavigatorMode _navigatorMode =
      _assignment?.navigatorMode ?? NavigatorMode.course;
  late bool _navigatorSummary = _assignment?.navigatorSummary ?? false;
  late DefenseFormat _defense = _assignment?.defense ?? DefenseFormat.none;
  late TeamRules? _team = _assignment?.team;
  late final _minSources = TextEditingController(
    text: '${_assignment?.searchMap?.minSources ?? 5}',
  );
  late final _minScholarly = TextEditingController(
    text: '${_assignment?.searchMap?.minScholarly ?? 50}',
  );
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    _title.dispose();
    _description.dispose();
    _maxPoints.dispose();
    _minSources.dispose();
    _minScholarly.dispose();
    super.dispose();
  }

  Future<void> _pickDue() async {
    final now = DateTime.now();
    final initial = _dueAt ?? now.add(const Duration(days: 7));
    final date = await showDatePicker(
      context: context,
      initialDate: initial,
      firstDate: DateTime(now.year - 1),
      lastDate: DateTime(now.year + 5),
    );
    if (date == null || !mounted) {
      return;
    }
    final time = await showTimePicker(
      context: context,
      initialTime: TimeOfDay.fromDateTime(_dueAt ?? DateTime(0, 1, 1, 23, 59)),
    );
    if (time == null) {
      return;
    }
    setState(
      () => _dueAt = DateTime(
        date.year,
        date.month,
        date.day,
        time.hour,
        time.minute,
      ),
    );
  }

  Future<void> _run(
    Future<void> Function() action, {
    bool deleted = false,
  }) async {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await action();
      ref.invalidate(assignmentsProvider(widget.courseId));
      final assignment = _assignment;
      if (assignment != null) {
        ref.invalidate(assignmentProvider(assignment.id));
      }
      router.pop(deleted);
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

  Future<void> _save() async {
    final l10n = context.l10n;
    final maxPoints = double.tryParse(
      _maxPoints.text.trim().replaceAll(',', '.'),
    );
    final minSources = int.tryParse(_minSources.text.trim());
    final minScholarly = int.tryParse(_minScholarly.text.trim());
    final mapInvalid =
        _mapRequired &&
        (minSources == null ||
            minSources < 0 ||
            minSources > 200 ||
            minScholarly == null ||
            minScholarly < 0 ||
            minScholarly > 100);
    if (_title.text.trim().isEmpty ||
        mapInvalid ||
        !(_team?.valid ?? true) ||
        (_grading == Grading.points && !(_weights?.valid ?? true)) ||
        (_grading == Grading.points && (maxPoints == null || maxPoints <= 0))) {
      setState(() => _error = l10n.errorInvalidAssignment);
      return;
    }
    final repository = ref.read(studyRepositoryProvider);
    final input = AssignmentInput(
      title: _title.text.trim(),
      description: _description.text.trim(),
      dueAt: _dueAt,
      grading: _grading,
      maxPoints: maxPoints,
      gradeWeights: _weights,
      published: _published,
      files: [...?widget.details?.files.map((file) => file.id)],
      searchMap: _mapRequired
          ? MapRules(minSources: minSources!, minScholarly: minScholarly!)
          : null,
      navigatorMode: _navigatorMode,
      navigatorSummary: _navigatorSummary,
      defense: _defense,
      team: _team,
    );
    final assignment = _assignment;
    await _run(
      () => assignment == null
          ? repository.create(widget.courseId, input)
          : repository.update(assignment.id, input),
    );
  }

  Future<void> _delete(String id, String title) async {
    final l10n = context.l10n;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.assignmentDeleteConfirm(title)),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: Text(l10n.actionCancel),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: Text(l10n.assignmentDelete),
          ),
        ],
      ),
    );
    if (confirmed ?? false) {
      await _run(
        () => ref.read(studyRepositoryProvider).delete(id),
        deleted: true,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final error = _error;
    final assignment = _assignment;
    return ListView(
      padding: const EdgeInsets.all(Spacing.l),
      children: [
        TextField(
          controller: _title,
          maxLength: 200,
          decoration: InputDecoration(labelText: l10n.assignmentTitle),
        ),
        const SizedBox(height: Spacing.m),
        TextField(
          controller: _description,
          minLines: 4,
          maxLines: 12,
          decoration: InputDecoration(
            labelText: l10n.assignmentDescription,
            alignLabelWithHint: true,
          ),
        ),
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.assignmentDueLabel,
          child: Wrap(
            spacing: Spacing.s,
            runSpacing: Spacing.s,
            crossAxisAlignment: WrapCrossAlignment.center,
            children: [
              LodestarButton(
                label: _dueAt == null
                    ? l10n.assignmentPickDue
                    : labels.time(_dueAt!),
                icon: Icons.event_outlined,
                variant: LodestarButtonVariant.secondary,
                onPressed: _pickDue,
              ),
              if (_dueAt != null)
                TextButton(
                  onPressed: () => setState(() => _dueAt = null),
                  child: Text(l10n.assignmentClearDue),
                ),
            ],
          ),
        ),
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.assignmentGrading,
          child: ChoiceGroup<Grading>(
            values: Grading.values,
            selected: _grading,
            labelOf: labels.grading,
            onSelected: (value) => setState(() => _grading = value),
          ),
        ),
        if (_grading == Grading.points) ...[
          const SizedBox(height: Spacing.l),
          TextField(
            controller: _maxPoints,
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
            decoration: InputDecoration(labelText: l10n.assignmentMaxPoints),
          ),
          const SizedBox(height: Spacing.s),
          GradeWeightsEditor(
            initial: _weights,
            onChanged: (value) => setState(() => _weights = value),
          ),
        ],
        const SizedBox(height: Spacing.l),
        SwitchListTile(
          contentPadding: EdgeInsets.zero,
          title: Text(l10n.assignmentMapRequired),
          value: _mapRequired,
          onChanged: (value) => setState(() => _mapRequired = value),
        ),
        if (_mapRequired)
          Row(
            children: [
              Expanded(
                child: TextField(
                  controller: _minSources,
                  keyboardType: TextInputType.number,
                  decoration: InputDecoration(
                    labelText: l10n.assignmentMapMinSources,
                  ),
                ),
              ),
              const SizedBox(width: Spacing.m),
              Expanded(
                child: TextField(
                  controller: _minScholarly,
                  keyboardType: TextInputType.number,
                  decoration: InputDecoration(
                    labelText: l10n.assignmentMapMinScholarly,
                  ),
                ),
              ),
            ],
          ),
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.assignmentNavigator,
          child: ChoiceGroup<NavigatorMode>(
            values: NavigatorMode.values,
            selected: _navigatorMode,
            labelOf: labels.navigatorMode,
            onSelected: (value) => setState(() => _navigatorMode = value),
          ),
        ),
        SwitchListTile(
          contentPadding: EdgeInsets.zero,
          title: Text(l10n.assignmentNavigatorSummary),
          subtitle: Text(l10n.assignmentNavigatorSummaryHint),
          value: _navigatorSummary,
          onChanged: (value) => setState(() => _navigatorSummary = value),
        ),
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.assignmentDefense,
          child: ChoiceGroup<DefenseFormat>(
            values: DefenseFormat.values,
            selected: _defense,
            labelOf: labels.defenseFormat,
            onSelected: (value) => setState(() => _defense = value),
          ),
        ),
        const SizedBox(height: Spacing.l),
        TeamRulesEditor(
          initial: _team,
          onChanged: (value) => setState(() => _team = value),
        ),
        const SizedBox(height: Spacing.l),
        SwitchListTile(
          contentPadding: EdgeInsets.zero,
          title: Text(l10n.assignmentPublish),
          subtitle: Text(l10n.assignmentPublishHint),
          value: _published,
          onChanged: (value) => setState(() => _published = value),
        ),
        if (error != null) FormError(message: error),
        const SizedBox(height: Spacing.xl),
        LodestarButton(label: l10n.actionSave, onPressed: _busy ? null : _save),
        if (assignment != null) ...[
          const SizedBox(height: Spacing.m),
          LodestarButton(
            label: l10n.assignmentDelete,
            icon: Icons.delete_outline,
            variant: LodestarButtonVariant.secondary,
            onPressed: _busy
                ? null
                : () => _delete(assignment.id, assignment.title),
          ),
        ],
      ],
    );
  }
}
