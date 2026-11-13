import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../study/course_person.dart';
import 'team_action.dart';
import 'team_details.dart';
import 'team_person_picker.dart';
import 'team_providers.dart';
import 'team_refresh.dart';

class TeamInviteSection extends ConsumerWidget {
  const TeamInviteSection({required this.details, super.key});

  final TeamDetails details;

  Future<void> _copy(BuildContext context, String code) async {
    final messenger = ScaffoldMessenger.of(context);
    final copied = context.l10n.teamCodeCopied;
    await Clipboard.setData(ClipboardData(text: code));
    messenger.showSnackBar(SnackBar(content: Text(copied)));
  }

  Future<void> _invite(BuildContext context, WidgetRef ref) async {
    final l10n = context.l10n;
    final scope = details.scope;
    final view = await ref.read(
      teamScopeProvider((scope.kind, scope.id)).future,
    );
    if (!context.mounted) {
      return;
    }
    final invited = {for (final person in details.invited) person.userId};
    final people = <CoursePerson>[
      for (final seeker in view.seekers)
        if (!invited.contains(seeker.person.userId)) seeker.person,
      for (final person in view.unplaced)
        if (!invited.contains(person.userId) &&
            !view.seekers.any((s) => s.person.userId == person.userId))
          person,
    ];
    final person = await TeamPersonPicker.show(
      context,
      title: l10n.teamInvite,
      people: people,
    );
    if (person == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).invite(details.team.id, person.userId),
    );
    ref.refreshTeam(details);
  }

  Future<void> _uninvite(
    BuildContext context,
    WidgetRef ref,
    CoursePerson person,
  ) async {
    await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).uninvite(details.team.id, person.userId),
    );
    ref.refreshTeam(details);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final code = details.inviteCode;
    final open =
        !details.team.archived &&
        !details.scope.archived &&
        details.team.members.length < details.scope.maxSize &&
        (details.manage || details.scope.selfFormed);
    if (!open && details.invited.isEmpty) {
      return const SizedBox.shrink();
    }
    return Section(
      title: l10n.teamInvitationsSent,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          if (open && code != null)
            ListTile(
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.link),
              title: SelectableText(code),
              subtitle: Text(
                l10n.teamCodeHint,
                style: TextStyle(color: context.palette.muted),
              ),
              trailing: IconButton(
                key: const ValueKey('team-copy-code'),
                tooltip: l10n.teamCopyCode,
                icon: const Icon(Icons.copy_outlined),
                onPressed: () => _copy(context, code),
              ),
            ),
          for (final person in details.invited)
            ListTile(
              key: ValueKey('invited-${person.userId}'),
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.hourglass_empty),
              title: Text(person.label),
              trailing: IconButton(
                tooltip: l10n.teamUninvite,
                icon: const Icon(Icons.close),
                onPressed: () => _uninvite(context, ref, person),
              ),
            ),
          if (open) ...[
            const SizedBox(height: Spacing.s),
            Align(
              alignment: Alignment.centerLeft,
              child: TextButton.icon(
                key: const ValueKey('team-invite'),
                onPressed: () => _invite(context, ref),
                icon: const Icon(Icons.person_add_alt_outlined),
                label: Text(l10n.teamInvite),
              ),
            ),
          ],
        ],
      ),
    );
  }
}
