import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'suggested_team.dart';
import 'team_action.dart';
import 'team_providers.dart';
import 'team_suggestion.dart';

class TeamSuggestionSection extends ConsumerStatefulWidget {
  const TeamSuggestionSection({
    required this.kind,
    required this.scopeId,
    this.onCreate,
    super.key,
  });

  final String kind;
  final String scopeId;
  final Future<void> Function(SuggestedTeam team)? onCreate;

  @override
  ConsumerState<TeamSuggestionSection> createState() =>
      _TeamSuggestionSectionState();
}

class _TeamSuggestionSectionState extends ConsumerState<TeamSuggestionSection> {
  TeamSuggestion? _suggestion;
  bool _busy = false;

  Future<void> _ask() async {
    setState(() => _busy = true);
    await TeamAction.run(context, () async {
      final suggestion = await ref
          .read(teamApiProvider)
          .suggest(widget.kind, widget.scopeId);
      setState(() => _suggestion = suggestion);
    });
    if (mounted) {
      setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    final suggestion = _suggestion;
    final create = widget.onCreate;
    return Section(
      title: l10n.teamSuggestTitle,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(l10n.teamSuggestHint, style: muted),
          const SizedBox(height: Spacing.s),
          Align(
            alignment: Alignment.centerLeft,
            child: OutlinedButton.icon(
              key: const ValueKey('team-suggest'),
              onPressed: _busy ? null : _ask,
              icon: const Icon(Icons.auto_awesome_outlined),
              label: Text(l10n.teamSuggest),
            ),
          ),
          if (suggestion != null) ...[
            const SizedBox(height: Spacing.m),
            if (suggestion.teams.isEmpty)
              Text(l10n.teamSuggestNone, style: muted),
            for (final team in suggestion.teams)
              LodestarCard(
                child: ListTile(
                  contentPadding: EdgeInsets.zero,
                  title: Text(team.members.map((p) => p.label).join(', ')),
                  subtitle: Text(team.why),
                  trailing: create == null
                      ? null
                      : TextButton(
                          onPressed: () => create(team),
                          child: Text(l10n.teamCreate),
                        ),
                ),
              ),
            if (suggestion.without.isNotEmpty)
              Text(
                l10n.teamSuggestWithout(
                  suggestion.without.map((p) => p.label).join(', '),
                ),
                style: muted,
              ),
          ],
        ],
      ),
    );
  }
}
