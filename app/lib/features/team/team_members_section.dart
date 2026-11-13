import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import '../study/course_person.dart';
import 'team_action.dart';
import 'team_details.dart';
import 'team_member.dart';
import 'team_member_action.dart';
import 'team_person_picker.dart';
import 'team_providers.dart';
import 'team_refresh.dart';
import 'team_text_dialog.dart';

class TeamMembersSection extends ConsumerWidget {
  const TeamMembersSection({
    required this.details,
    required this.me,
    super.key,
  });

  final TeamDetails details;
  final String? me;

  Future<void> _label(
    BuildContext context,
    WidgetRef ref,
    TeamMember member,
  ) async {
    final l10n = context.l10n;
    final label = await TeamTextDialog.show(
      context,
      title: l10n.teamMyLabel,
      label: l10n.teamLabel,
      hint: l10n.teamLabelHint,
      action: l10n.actionSave,
      initial: member.label,
      required: false,
    );
    if (label == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).label(details.team.id, label),
    );
    ref.refreshTeam(details);
  }

  Future<void> _add(BuildContext context, WidgetRef ref) async {
    final l10n = context.l10n;
    final scope = details.scope;
    final view = await ref.read(
      teamScopeProvider((scope.kind, scope.id)).future,
    );
    if (!context.mounted) {
      return;
    }
    final people = <CoursePerson>[
      ...view.unplaced,
      for (final team in view.teams)
        if (team.id != details.team.id)
          for (final member in team.members) member.person,
    ];
    final person = await TeamPersonPicker.show(
      context,
      title: l10n.teamAddMember,
      people: people,
    );
    if (person == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).add(details.team.id, person.userId),
    );
    ref.refreshTeam(details);
  }

  Future<void> _manage(
    BuildContext context,
    WidgetRef ref,
    TeamMember member,
    TeamMemberAction action,
  ) async {
    await TeamAction.run(
      context,
      () => ref
          .read(teamApiProvider)
          .remove(
            details.team.id,
            member.person.userId,
            solo: action == TeamMemberAction.separate,
          ),
    );
    ref.refreshTeam(details);
  }

  Widget? _trailing(BuildContext context, WidgetRef ref, TeamMember member) {
    final l10n = context.l10n;
    if (details.team.archived) {
      return null;
    }
    if (member.person.userId == me) {
      return IconButton(
        key: const ValueKey('team-label'),
        tooltip: l10n.teamMyLabel,
        icon: const Icon(Icons.sell_outlined),
        onPressed: () => _label(context, ref, member),
      );
    }
    if (!details.manage) {
      return null;
    }
    return PopupMenuButton<TeamMemberAction>(
      key: ValueKey('member-menu-${member.person.userId}'),
      onSelected: (action) => _manage(context, ref, member, action),
      itemBuilder: (context) => [
        PopupMenuItem(
          value: TeamMemberAction.remove,
          child: Text(l10n.teamRemoveMember),
        ),
        PopupMenuItem(
          value: TeamMemberAction.separate,
          child: Text(l10n.teamSeparateMember),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return Section(
      title: l10n.teamMembers(details.team.members.length),
      child: Column(
        children: [
          for (final member in details.team.members)
            ListTile(
              key: ValueKey('member-${member.person.userId}'),
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.person_outline),
              title: Text(member.person.label),
              subtitle: member.label.isEmpty ? null : Text(member.label),
              trailing: _trailing(context, ref, member),
            ),
          if (details.manage && !details.team.archived)
            Align(
              alignment: Alignment.centerLeft,
              child: TextButton.icon(
                key: const ValueKey('team-add'),
                onPressed: () => _add(context, ref),
                icon: const Icon(Icons.person_add_alt_outlined),
                label: Text(l10n.teamAddMember),
              ),
            ),
        ],
      ),
    );
  }
}
