import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import '../team/team_details.dart';
import '../team/team_providers.dart';
import '../team/team_refresh.dart';

class LeagueSolutionForm extends ConsumerStatefulWidget {
  const LeagueSolutionForm({required this.team, super.key});

  final TeamDetails team;

  @override
  ConsumerState<LeagueSolutionForm> createState() => _LeagueSolutionFormState();
}

class _LeagueSolutionFormState extends ConsumerState<LeagueSolutionForm> {
  final _text = TextEditingController();
  final _links = TextEditingController();
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    _text.dispose();
    _links.dispose();
    super.dispose();
  }

  Future<void> _propose() async {
    final l10n = context.l10n;
    final text = _text.text.trim();
    final links = [
      for (final line in _links.text.split('\n'))
        if (line.trim().isNotEmpty) line.trim(),
    ];
    if (text.isEmpty && links.isEmpty) {
      setState(() => _error = l10n.errorLeagueSolution);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    final done = await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).propose(widget.team.team.id, {
        'text': text,
        'links': links,
      }, ''),
    );
    ref.refreshTeam(widget.team);
    if (!mounted) {
      return;
    }
    setState(() => _busy = false);
    if (done) {
      _text.clear();
      _links.clear();
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(l10n.teamProposeHint),
        const SizedBox(height: Spacing.m),
        TextField(
          key: const ValueKey('solution-text'),
          controller: _text,
          minLines: 4,
          maxLines: 12,
          decoration: InputDecoration(
            labelText: l10n.leagueSolutionText,
            alignLabelWithHint: true,
          ),
        ),
        const SizedBox(height: Spacing.m),
        TextField(
          key: const ValueKey('solution-links'),
          controller: _links,
          minLines: 1,
          maxLines: 5,
          keyboardType: TextInputType.url,
          decoration: InputDecoration(labelText: l10n.leagueSolutionLinks),
        ),
        if (error != null) FormError(message: error),
        const SizedBox(height: Spacing.l),
        Align(
          alignment: Alignment.centerLeft,
          child: LodestarButton(
            key: const ValueKey('solution-propose'),
            label: widget.team.proposal == null
                ? l10n.teamPropose
                : l10n.teamProposeAgain,
            icon: Icons.how_to_vote_outlined,
            onPressed: _busy ? null : _propose,
          ),
        ),
      ],
    );
  }
}
