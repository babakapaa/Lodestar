import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/navigator_tile.dart';
import '../study/course_person.dart';
import 'suggested_team.dart';
import 'team_action.dart';
import 'team_providers.dart';
import 'team_scope_refresh.dart';
import 'team_scope_view.dart';
import 'team_seekers_section.dart';
import 'team_suggestion_section.dart';
import 'team_text_dialog.dart';

class TeamScopeStudent extends ConsumerWidget {
  const TeamScopeStudent({required this.view, super.key});

  final TeamScopeView view;

  Future<void> _seek(BuildContext context, WidgetRef ref) async {
    final l10n = context.l10n;
    final scope = view.scope;
    if (view.seeking) {
      await TeamAction.run(
        context,
        () => ref
            .read(teamApiProvider)
            .seek(scope.kind, scope.id, seeking: false),
      );
      ref.refreshScope(scope);
      return;
    }
    final note = await TeamTextDialog.show(
      context,
      title: l10n.teamSeeking,
      label: l10n.teamSeekingNote,
      action: l10n.teamSeekingStart,
      required: false,
      lines: 2,
    );
    if (note == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref
          .read(teamApiProvider)
          .seek(scope.kind, scope.id, seeking: true, note: note),
    );
    ref.refreshScope(scope);
  }

  Future<void> _create(
    BuildContext context,
    WidgetRef ref, [
    SuggestedTeam? suggested,
  ]) async {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    final scope = view.scope;
    final name = await TeamTextDialog.show(
      context,
      title: l10n.teamCreate,
      label: l10n.teamName,
      action: l10n.teamCreate,
    );
    if (name == null || !context.mounted) {
      return;
    }
    String? id;
    final done = await TeamAction.run(context, () async {
      final api = ref.read(teamApiProvider);
      final details = await api.create(scope.kind, scope.id, name);
      id = details.team.id;
      for (final person in suggested?.members ?? const <CoursePerson>[]) {
        if (!details.team.has(person.userId)) {
          await api.invite(details.team.id, person.userId);
        }
      }
    });
    ref.refreshScope(scope);
    final created = id;
    if (done && created != null) {
      await router.push(AppRoutes.team(created));
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final scope = view.scope;
    final team = view.myTeam;
    if (team != null) {
      return NavigatorTile(
        key: const ValueKey('my-team'),
        icon: Icons.groups_outlined,
        title: team.name,
        subtitle: team.names,
        route: AppRoutes.team(team.id),
      );
    }
    if (view.solo) {
      return NoticeBanner(icon: Icons.person_outline, message: l10n.teamSolo);
    }
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(scope.selfFormed ? l10n.teamNoneSelf : l10n.teamNoneByTeacher),
        const SizedBox(height: Spacing.l),
        Wrap(
          spacing: Spacing.s,
          runSpacing: Spacing.s,
          children: [
            if (scope.selfFormed && !scope.archived)
              LodestarButton(
                key: const ValueKey('team-create'),
                label: l10n.teamCreate,
                icon: Icons.group_add_outlined,
                onPressed: () => _create(context, ref),
              ),
            LodestarButton(
              key: const ValueKey('team-seek'),
              label: view.seeking
                  ? l10n.teamSeekingStop
                  : l10n.teamSeekingStart,
              icon: Icons.person_search_outlined,
              variant: LodestarButtonVariant.secondary,
              onPressed: scope.archived ? null : () => _seek(context, ref),
            ),
          ],
        ),
        if (view.seekers.isNotEmpty) ...[
          const SizedBox(height: Spacing.xl),
          TeamSeekersSection(seekers: view.seekers),
        ],
        if (scope.selfFormed && !scope.archived) ...[
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
