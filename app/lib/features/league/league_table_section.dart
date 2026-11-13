import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'league_labels_context.dart';
import 'league_table_row.dart';

class LeagueTableSection extends StatelessWidget {
  const LeagueTableSection({
    required this.division,
    required this.rows,
    super.key,
  });

  final String division;
  final List<LeagueTableRow> rows;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.leagueLabels;
    return Section(
      title: l10n.leagueTable(division),
      child: Column(
        children: [
          for (final row in rows)
            ListTile(
              contentPadding: EdgeInsets.zero,
              selected: row.me,
              leading: Text('${row.place}'),
              title: Text(row.me ? l10n.leagueYou(row.name) : row.name),
              subtitle: Text(
                [
                  labels.growth(row.growth),
                  if (row.topGrowth) l10n.leagueTopGrowth,
                ].join(', '),
              ),
              trailing: Text('${row.points}'),
            ),
        ],
      ),
    );
  }
}
