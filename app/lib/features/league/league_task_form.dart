import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/study_labels_context.dart';
import 'league_providers.dart';
import 'league_task.dart';
import 'task_input.dart';

class LeagueTaskForm extends ConsumerStatefulWidget {
  const LeagueTaskForm({
    required this.leagueId,
    required this.onSaved,
    this.task,
    super.key,
  });

  final String leagueId;
  final LeagueTask? task;
  final VoidCallback onSaved;

  @override
  ConsumerState<LeagueTaskForm> createState() => _LeagueTaskFormState();
}

class _LeagueTaskFormState extends ConsumerState<LeagueTaskForm> {
  late final _title = TextEditingController(text: widget.task?.title ?? '');
  late final _about = TextEditingController(text: widget.task?.about ?? '');
  late final _partner = TextEditingController(text: widget.task?.partner ?? '');
  late final _points = TextEditingController(
    text: '${widget.task?.points ?? 50}',
  );
  late final _minSize = TextEditingController(
    text: '${widget.task?.minSize ?? 2}',
  );
  late final _maxSize = TextEditingController(
    text: '${widget.task?.maxSize ?? 4}',
  );
  late bool _selfFormed = widget.task?.selfFormed ?? true;
  late bool _closed = widget.task?.closed ?? false;
  late DateTime? _dueAt = widget.task?.dueAt?.toLocal();
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    for (final c in [_title, _about, _partner, _points, _minSize, _maxSize]) {
      c.dispose();
    }
    super.dispose();
  }

  Future<void> _pickDue() async {
    final now = DateTime.now();
    final date = await showDatePicker(
      context: context,
      initialDate: _dueAt ?? now.add(const Duration(days: 14)),
      firstDate: DateTime(now.year - 1),
      lastDate: DateTime(now.year + 5),
    );
    if (date != null) {
      setState(
        () => _dueAt = DateTime(date.year, date.month, date.day, 23, 59),
      );
    }
  }

  Future<void> _save() async {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    final points = int.tryParse(_points.text.trim());
    final minSize = int.tryParse(_minSize.text.trim());
    final maxSize = int.tryParse(_maxSize.text.trim());
    if (_title.text.trim().isEmpty ||
        points == null ||
        points < 1 ||
        minSize == null ||
        maxSize == null ||
        minSize < 2 ||
        maxSize < minSize ||
        maxSize > 12) {
      setState(() => _error = l10n.errorLeagueTask);
      return;
    }
    final input = TaskInput(
      title: _title.text.trim(),
      about: _about.text.trim(),
      partner: _partner.text.trim(),
      points: points,
      minSize: minSize,
      maxSize: maxSize,
      selfFormed: _selfFormed,
      dueAt: _dueAt,
      closed: _closed,
    );
    setState(() {
      _busy = true;
      _error = null;
    });
    final api = ref.read(leagueApiProvider);
    final task = widget.task;
    try {
      if (task == null) {
        await api.createTask(widget.leagueId, input);
      } else {
        await api.updateTask(task.id, input);
      }
      widget.onSaved();
      router.pop();
    } on ApiException catch (error) {
      setState(() => _error = ApiErrorText.of(l10n, error));
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    final due = _dueAt;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        TextField(
          key: const ValueKey('task-title'),
          controller: _title,
          maxLength: 200,
          decoration: InputDecoration(labelText: l10n.leagueTaskTitle),
        ),
        TextField(
          controller: _about,
          minLines: 3,
          maxLines: 10,
          decoration: InputDecoration(
            labelText: l10n.leagueTaskAbout,
            alignLabelWithHint: true,
          ),
        ),
        const SizedBox(height: Spacing.m),
        TextField(
          controller: _partner,
          decoration: InputDecoration(
            labelText: l10n.leagueTaskPartner,
            helperText: l10n.leagueTaskPartnerHint,
          ),
        ),
        const SizedBox(height: Spacing.m),
        TextField(
          key: const ValueKey('task-points'),
          controller: _points,
          keyboardType: TextInputType.number,
          decoration: InputDecoration(labelText: l10n.leagueTaskPoints),
        ),
        const SizedBox(height: Spacing.m),
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: _minSize,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(labelText: l10n.teamMinSize),
              ),
            ),
            const SizedBox(width: Spacing.m),
            Expanded(
              child: TextField(
                controller: _maxSize,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(labelText: l10n.teamMaxSize),
              ),
            ),
          ],
        ),
        SwitchListTile(
          contentPadding: EdgeInsets.zero,
          title: Text(l10n.teamSelfFormedSwitch),
          value: _selfFormed,
          onChanged: (value) => setState(() => _selfFormed = value),
        ),
        Align(
          alignment: Alignment.centerLeft,
          child: LodestarButton(
            label: due == null
                ? l10n.assignmentPickDue
                : context.studyLabels.time(due),
            icon: Icons.event_outlined,
            variant: LodestarButtonVariant.secondary,
            onPressed: _pickDue,
          ),
        ),
        if (widget.task != null)
          SwitchListTile(
            contentPadding: EdgeInsets.zero,
            title: Text(l10n.leagueTaskClose),
            value: _closed,
            onChanged: (value) => setState(() => _closed = value),
          ),
        if (error != null) FormError(message: error),
        const SizedBox(height: Spacing.xl),
        LodestarButton(
          key: const ValueKey('task-save'),
          label: l10n.actionSave,
          onPressed: _busy ? null : _save,
        ),
      ],
    );
  }
}
