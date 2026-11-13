import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'work_weight.dart';

class LeagueFormulaSection extends StatelessWidget {
  const LeagueFormulaSection({required this.formula, super.key});

  final List<WorkWeight> formula;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return Section(
      title: l10n.leagueFormula,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(l10n.leagueFormulaHint),
          for (final weight in formula)
            ListTile(
              dense: true,
              contentPadding: EdgeInsets.zero,
              title: Text(weight.title),
              trailing: Text(l10n.leagueWeight(weight.weight)),
            ),
        ],
      ),
    );
  }
}
