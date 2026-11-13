import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'team_seeker.dart';

class TeamSeekersSection extends StatelessWidget {
  const TeamSeekersSection({required this.seekers, this.trailing, super.key});

  final List<TeamSeeker> seekers;
  final Widget Function(TeamSeeker seeker)? trailing;

  @override
  Widget build(BuildContext context) {
    final action = trailing;
    return Section(
      title: context.l10n.teamSeekers,
      child: Column(
        children: [
          for (final seeker in seekers)
            ListTile(
              key: ValueKey('seeker-${seeker.person.userId}'),
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.person_search_outlined),
              title: Text(seeker.person.label),
              subtitle: seeker.note.isEmpty ? null : Text(seeker.note),
              trailing: action?.call(seeker),
            ),
        ],
      ),
    );
  }
}
