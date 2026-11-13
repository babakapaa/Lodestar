import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'history_entry.dart';
import 'team_labels_context.dart';

class TeamHistorySection extends StatelessWidget {
  const TeamHistorySection({required this.history, super.key});

  final List<HistoryEntry> history;

  @override
  Widget build(BuildContext context) {
    final labels = context.teamLabels;
    return Section(
      title: context.l10n.teamHistory,
      child: Column(
        children: [
          for (final entry in history)
            ListTile(
              dense: true,
              contentPadding: EdgeInsets.zero,
              title: Text(labels.event(entry)),
              subtitle: Text(labels.time(entry.at)),
            ),
        ],
      ),
    );
  }
}
