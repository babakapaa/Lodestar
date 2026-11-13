import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/assignment_details.dart';
import '../study/study_providers.dart';
import 'team_action.dart';
import 'team_details.dart';
import 'team_providers.dart';
import 'team_refresh.dart';

class TeamProposeForm extends ConsumerStatefulWidget {
  const TeamProposeForm({
    required this.assignment,
    required this.team,
    super.key,
  });

  final AssignmentDetails assignment;
  final TeamDetails team;

  @override
  ConsumerState<TeamProposeForm> createState() => _TeamProposeFormState();
}

class _TeamProposeFormState extends ConsumerState<TeamProposeForm> {
  final _body = TextEditingController();
  final _note = TextEditingController();
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    _body.dispose();
    _note.dispose();
    super.dispose();
  }

  Future<void> _propose() async {
    final l10n = context.l10n;
    final body = _body.text.trim();
    if (body.isEmpty) {
      setState(() => _error = l10n.errorInvalidSubmission);
      return;
    }
    final assignment = widget.assignment.assignment;
    final map = ref.read(searchMapProvider(assignment.id)).value?.value.map;
    final attach =
        map != null && (assignment.searchMap != null || !map.isEmpty);
    setState(() {
      _busy = true;
      _error = null;
    });
    final done = await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).propose(widget.team.team.id, {
        'body': body,
        if (attach) 'searchMap': map.toJson(),
      }, _note.text.trim()),
    );
    ref.refreshTeam(widget.team);
    if (!mounted) {
      return;
    }
    setState(() => _busy = false);
    if (done) {
      _body.clear();
      _note.clear();
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    final assignment = widget.assignment.assignment;
    final rules = assignment.searchMap;
    final map = ref.watch(searchMapProvider(assignment.id)).value?.value.map;
    final blocked =
        rules != null && (map == null || map.gaps(rules).isNotEmpty);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(l10n.teamProposeHint),
        const SizedBox(height: Spacing.m),
        TextField(
          key: const ValueKey('team-work'),
          controller: _body,
          minLines: 4,
          maxLines: 12,
          decoration: InputDecoration(
            labelText: l10n.submitHint,
            alignLabelWithHint: true,
          ),
        ),
        const SizedBox(height: Spacing.m),
        TextField(
          controller: _note,
          decoration: InputDecoration(labelText: l10n.teamProposeNote),
        ),
        if (error != null) FormError(message: error),
        if (blocked) ...[
          const SizedBox(height: Spacing.m),
          NoticeBanner(icon: Icons.travel_explore, message: l10n.mapNotReady),
        ],
        const SizedBox(height: Spacing.l),
        Align(
          alignment: Alignment.centerLeft,
          child: LodestarButton(
            key: const ValueKey('team-propose'),
            label: widget.team.proposal == null
                ? l10n.teamPropose
                : l10n.teamProposeAgain,
            icon: Icons.how_to_vote_outlined,
            onPressed: _busy || blocked ? null : _propose,
          ),
        ),
      ],
    );
  }
}
