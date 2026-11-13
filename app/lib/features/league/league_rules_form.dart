import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'division.dart';
import 'division_draft.dart';
import 'division_fields.dart';
import 'league_providers.dart';
import 'league_refresh.dart';
import 'league_view.dart';

class LeagueRulesForm extends ConsumerStatefulWidget {
  const LeagueRulesForm({required this.view, super.key});

  final LeagueView view;

  @override
  ConsumerState<LeagueRulesForm> createState() => _LeagueRulesFormState();
}

class _LeagueRulesFormState extends ConsumerState<LeagueRulesForm> {
  late final List<DivisionDraft> _divisions = [
    for (final d in widget.view.divisions) DivisionDraft(d),
  ];
  late final Map<String, TextEditingController> _weights = {
    for (final w in widget.view.formula)
      '${w.kind}/${w.id}': TextEditingController(text: '${w.weight}'),
  };
  late bool _confirm = widget.view.confirm;
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    for (final d in _divisions) {
      d.dispose();
    }
    for (final c in _weights.values) {
      c.dispose();
    }
    super.dispose();
  }

  void _add() {
    final last = _divisions.lastOrNull?.value;
    setState(
      () => _divisions.add(
        DivisionDraft(
          Division(name: '', from: (last?.from ?? 0) + 150, gpuHours: 0),
        ),
      ),
    );
  }

  void _remove(int index) {
    setState(() => _divisions.removeAt(index).dispose());
  }

  Future<void> _save() async {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    final divisions = [for (final d in _divisions) d.value];
    final weights = [
      for (final w in widget.view.formula)
        (w, int.tryParse(_weights['${w.kind}/${w.id}']!.text.trim())),
    ];
    if (divisions.contains(null) ||
        divisions.isEmpty ||
        weights.any((w) => w.$2 == null || w.$2! < 0)) {
      setState(() => _error = l10n.errorLeagueRules);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    final view = widget.view;
    final done = await TeamAction.run(
      context,
      () => ref
          .read(leagueApiProvider)
          .rules(
            view.id,
            divisions: [for (final d in divisions) d!],
            confirm: _confirm,
            weights: [for (final (w, value) in weights) w.toJson(value)],
          ),
    );
    ref.refreshLeague(view.kind, view.originId);
    if (!mounted) {
      return;
    }
    setState(() => _busy = false);
    if (done) {
      router.pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(
          l10n.leagueRules,
          style: Theme.of(context).textTheme.headlineSmall,
        ),
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.leagueDivisions,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(l10n.leagueDivisionsHint, style: muted),
              const SizedBox(height: Spacing.m),
              for (final (i, draft) in _divisions.indexed)
                DivisionFields(
                  key: ObjectKey(draft),
                  index: i,
                  draft: draft,
                  onRemove: _divisions.length > 1 ? () => _remove(i) : null,
                ),
              Align(
                alignment: Alignment.centerLeft,
                child: TextButton.icon(
                  key: const ValueKey('division-add'),
                  onPressed: _add,
                  icon: const Icon(Icons.add),
                  label: Text(l10n.leagueDivisionAdd),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: Spacing.l),
        SwitchListTile(
          key: const ValueKey('league-confirm'),
          contentPadding: EdgeInsets.zero,
          title: Text(l10n.leagueConfirm),
          subtitle: Text(l10n.leagueConfirmHint),
          value: _confirm,
          onChanged: (value) => setState(() => _confirm = value),
        ),
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.leagueFormula,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(l10n.leagueWeightsHint, style: muted),
              for (final w in widget.view.formula)
                TextField(
                  key: ValueKey('weight-${w.id}'),
                  controller: _weights['${w.kind}/${w.id}'],
                  keyboardType: TextInputType.number,
                  decoration: InputDecoration(labelText: w.title),
                ),
            ],
          ),
        ),
        if (error != null) FormError(message: error),
        const SizedBox(height: Spacing.xl),
        LodestarButton(
          key: const ValueKey('rules-save'),
          label: l10n.actionSave,
          onPressed: _busy ? null : _save,
        ),
      ],
    );
  }
}
