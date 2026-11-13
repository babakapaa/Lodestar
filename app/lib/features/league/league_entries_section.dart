import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'league_entry.dart';
import 'league_labels_context.dart';

class LeagueEntriesSection extends StatelessWidget {
  const LeagueEntriesSection({
    required this.entries,
    this.onDispute,
    super.key,
  });

  final List<LeagueEntry> entries;
  final void Function(LeagueEntry entry)? onDispute;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.leagueLabels;
    final dispute = onDispute;
    return Section(
      title: l10n.leagueEntries,
      child: Column(
        children: [
          if (entries.isEmpty) Text(l10n.leagueEntriesNone),
          for (final entry in entries)
            ListTile(
              key: ValueKey('entry-${entry.id}'),
              contentPadding: EdgeInsets.zero,
              title: Text(entry.title),
              subtitle: Text(labels.entry(entry)),
              leading: Text(labels.points(entry.points)),
              trailing: dispute == null
                  ? null
                  : IconButton(
                      key: ValueKey('dispute-${entry.id}'),
                      tooltip: l10n.leagueDispute,
                      icon: const Icon(Icons.report_outlined),
                      onPressed: () => dispute(entry),
                    ),
            ),
        ],
      ),
    );
  }
}
