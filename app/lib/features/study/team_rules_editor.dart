import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'study_labels_context.dart';
import 'team_rules.dart';

class TeamRulesEditor extends StatefulWidget {
  const TeamRulesEditor({
    required this.initial,
    required this.onChanged,
    super.key,
  });

  final TeamRules? initial;
  final ValueChanged<TeamRules?> onChanged;

  @override
  State<TeamRulesEditor> createState() => _TeamRulesEditorState();
}

class _TeamRulesEditorState extends State<TeamRulesEditor> {
  late bool _enabled = widget.initial != null;
  late bool _selfFormed = widget.initial?.selfFormed ?? true;
  late DateTime? _freeLeaveUntil = widget.initial?.freeLeaveUntil?.toLocal();
  late final _minSize = TextEditingController(
    text: '${widget.initial?.minSize ?? 2}',
  );
  late final _maxSize = TextEditingController(
    text: '${widget.initial?.maxSize ?? 4}',
  );

  @override
  void dispose() {
    _minSize.dispose();
    _maxSize.dispose();
    super.dispose();
  }

  void _emit() {
    if (!_enabled) {
      widget.onChanged(null);
      return;
    }
    widget.onChanged(
      TeamRules(
        minSize: int.tryParse(_minSize.text.trim()) ?? 0,
        maxSize: int.tryParse(_maxSize.text.trim()) ?? 0,
        selfFormed: _selfFormed,
        freeLeaveUntil: _freeLeaveUntil,
      ),
    );
  }

  void _update(VoidCallback change) {
    setState(change);
    _emit();
  }

  Future<void> _pickFreeLeave() async {
    final now = DateTime.now();
    final date = await showDatePicker(
      context: context,
      initialDate: _freeLeaveUntil ?? now.add(const Duration(days: 7)),
      firstDate: DateTime(now.year - 1),
      lastDate: DateTime(now.year + 5),
    );
    if (date != null) {
      _update(
        () =>
            _freeLeaveUntil = DateTime(date.year, date.month, date.day, 23, 59),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final until = _freeLeaveUntil;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SwitchListTile(
          key: const ValueKey('assignment-team'),
          contentPadding: EdgeInsets.zero,
          title: Text(l10n.assignmentTeam),
          subtitle: Text(l10n.assignmentTeamHint),
          value: _enabled,
          onChanged: (value) => _update(() => _enabled = value),
        ),
        if (_enabled) ...[
          Row(
            children: [
              Expanded(
                child: TextField(
                  key: const ValueKey('team-min'),
                  controller: _minSize,
                  keyboardType: TextInputType.number,
                  decoration: InputDecoration(labelText: l10n.teamMinSize),
                  onChanged: (_) => _emit(),
                ),
              ),
              const SizedBox(width: Spacing.m),
              Expanded(
                child: TextField(
                  key: const ValueKey('team-max'),
                  controller: _maxSize,
                  keyboardType: TextInputType.number,
                  decoration: InputDecoration(labelText: l10n.teamMaxSize),
                  onChanged: (_) => _emit(),
                ),
              ),
            ],
          ),
          SwitchListTile(
            contentPadding: EdgeInsets.zero,
            title: Text(l10n.teamSelfFormedSwitch),
            value: _selfFormed,
            onChanged: (value) => _update(() => _selfFormed = value),
          ),
          Wrap(
            spacing: Spacing.s,
            runSpacing: Spacing.s,
            crossAxisAlignment: WrapCrossAlignment.center,
            children: [
              LodestarButton(
                label: until == null
                    ? l10n.teamPickFreeLeave
                    : l10n.teamFreeLeaveUntil(context.studyLabels.time(until)),
                icon: Icons.event_outlined,
                variant: LodestarButtonVariant.secondary,
                onPressed: _pickFreeLeave,
              ),
              if (until != null)
                TextButton(
                  onPressed: () => _update(() => _freeLeaveUntil = null),
                  child: Text(l10n.assignmentClearDue),
                ),
            ],
          ),
        ],
      ],
    );
  }
}
