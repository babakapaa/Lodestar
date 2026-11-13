import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'league_dispute.dart';
import 'league_entry.dart';
import 'league_labels_context.dart';

class LeagueDisputesSection extends StatelessWidget {
  const LeagueDisputesSection({
    required this.disputes,
    this.entries = const [],
    this.onAnswer,
    super.key,
  });

  final List<LeagueDispute> disputes;
  final List<LeagueEntry> entries;
  final void Function(LeagueDispute dispute)? onAnswer;

  String? _entryTitle(LeagueDispute dispute) {
    if (dispute.entry case final entry?) {
      return entry.title;
    }
    for (final entry in entries) {
      if (entry.id == dispute.entryId) {
        return entry.title;
      }
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    final labels = context.leagueLabels;
    final answer = onAnswer;
    return Section(
      title: context.l10n.leagueDisputes,
      child: Column(
        children: [
          for (final dispute in disputes)
            ListTile(
              key: ValueKey('league-dispute-${dispute.id}'),
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.report_outlined),
              title: Text(
                [
                  if (dispute.person case final person?) person.label,
                  ?_entryTitle(dispute),
                  dispute.text,
                ].join(': '),
              ),
              subtitle: Text(
                [
                  labels.disputeStatus(dispute.status),
                  if (dispute.answer.isNotEmpty) dispute.answer,
                ].join(': '),
              ),
              onTap: answer == null ? null : () => answer(dispute),
            ),
        ],
      ),
    );
  }
}
