import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/course_person.dart';
import '../study/navigator_tile.dart';
import 'leave_request_tile.dart';
import 'suggested_team.dart';
import 'team_action.dart';
import 'team_providers.dart';
import 'team_scope_refresh.dart';
import 'team_scope_view.dart';
import 'team_seekers_section.dart';
import 'team_suggestion_section.dart';
import 'team_text_dialog.dart';

class TeamScopeLead extends ConsumerWidget {
  const TeamScopeLead({required this.view, super.key});

  final TeamScopeView view;

  Future<void> _create(
    BuildContext context,
    WidgetRef ref, [
    SuggestedTeam? suggested,
  ]) async {
    final l10n = context.l10n;
    final scope = view.scope;
    final members = suggested?.members ?? const <CoursePerson>[];
    final name = await TeamTextDialog.show(
      context,
      title: l10n.teamCreate,
      label: l10n.teamName,
      action: l10n.teamCreate,
      initial: members.map((p) => p.label.split(' ').first).join(', '),
    );
    if (name == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).create(scope.kind, scope.id, name, [
        for (final person in members) person.userId,
      ]),
    );
    ref.refreshScope(scope);
  }

  Future<void> _distribute(BuildContext context, WidgetRef ref) async {
    final scope = view.scope;
    await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).distribute(scope.kind, scope.id),
    );
    ref.refreshScope(scope);
  }

  Widget _people(String title, List<CoursePerson> people, IconData icon) =>
      Section(
        title: title,
        child: Column(
          children: [
            for (final person in people)
              ListTile(
                contentPadding: EdgeInsets.zero,
                leading: Icon(icon),
                title: Text(person.label),
              ),
          ],
        ),
      );

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final scope = view.scope;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (!scope.archived)
          Wrap(
            spacing: Spacing.s,
            runSpacing: Spacing.s,
            children: [
              LodestarButton(
                key: const ValueKey('team-create'),
                label: l10n.teamCreate,
                icon: Icons.group_add_outlined,
                onPressed: () => _create(context, ref),
              ),
              LodestarButton(
                key: const ValueKey('team-distribute'),
                label: l10n.teamDistribute,
                icon: Icons.shuffle,
                variant: LodestarButtonVariant.secondary,
                onPressed: view.unplaced.isEmpty
                    ? null
                    : () => _distribute(context, ref),
              ),
            ],
          ),
        if (view.requests.isNotEmpty) ...[
          const SizedBox(height: Spacing.xl),
          Section(
            title: l10n.teamLeaveRequests,
            child: Column(
              children: [
                for (final request in view.requests)
                  LeaveRequestTile(
                    request: request,
                    onDecided: () => ref.refreshScope(scope),
                  ),
              ],
            ),
          ),
        ],
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.teamsAll(view.teams.length),
          child: Column(
            children: [
              if (view.teams.isEmpty) Text(l10n.teamsNone),
              for (final team in view.teams) ...[
                NavigatorTile(
                  key: ValueKey('team-${team.id}'),
                  icon: Icons.groups_outlined,
                  title: team.name,
                  subtitle: team.names,
                  route: AppRoutes.team(team.id),
                ),
                const SizedBox(height: Spacing.s),
              ],
            ],
          ),
        ),
        if (view.unplaced.isNotEmpty) ...[
          const SizedBox(height: Spacing.xl),
          _people(
            l10n.teamUnplaced(view.unplaced.length),
            view.unplaced,
            Icons.person_outline,
          ),
        ],
        if (view.seekers.isNotEmpty) ...[
          const SizedBox(height: Spacing.xl),
          TeamSeekersSection(seekers: view.seekers),
        ],
        if (view.solos.isNotEmpty) ...[
          const SizedBox(height: Spacing.xl),
          _people(l10n.teamSolos, view.solos, Icons.person_off_outlined),
        ],
        if (view.unplaced.length > 1 && !scope.archived) ...[
          const SizedBox(height: Spacing.xl),
          TeamSuggestionSection(
            kind: scope.kind,
            scopeId: scope.id,
            onCreate: (team) => _create(context, ref, team),
          ),
        ],
      ],
    );
  }
}
