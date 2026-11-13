import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'award_grant.dart';
import 'league_labels_context.dart';

class LeagueGrantsSection extends StatelessWidget {
  const LeagueGrantsSection({
    required this.grants,
    this.nameOf,
    this.onRevoke,
    super.key,
  });

  final List<AwardGrant> grants;
  final String Function(String userId)? nameOf;
  final void Function(AwardGrant grant)? onRevoke;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.leagueLabels;
    final name = nameOf;
    final revoke = onRevoke;
    return Section(
      title: l10n.leagueAwards,
      child: Column(
        children: [
          for (final grant in grants)
            ListTile(
              key: ValueKey('grant-${grant.id}'),
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.emoji_events_outlined),
              title: Text(
                [if (name != null) name(grant.userId), grant.title].join(': '),
              ),
              subtitle: Text(
                [
                  if (grant.gpuHours > 0) l10n.leagueGpuHours(grant.gpuHours),
                  if (grant.reason.isNotEmpty) grant.reason,
                  labels.time(grant.grantedAt),
                ].join(', '),
              ),
              trailing: revoke == null
                  ? null
                  : IconButton(
                      tooltip: l10n.leagueRevoke,
                      icon: const Icon(Icons.close),
                      onPressed: () => revoke(grant),
                    ),
            ),
        ],
      ),
    );
  }
}
