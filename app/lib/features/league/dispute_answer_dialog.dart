import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'dispute_answer.dart';
import 'league_dispute.dart';

class DisputeAnswerDialog extends StatefulWidget {
  const DisputeAnswerDialog({required this.dispute, super.key});

  final LeagueDispute dispute;

  static Future<DisputeAnswer?> show(
    BuildContext context,
    LeagueDispute dispute,
  ) => showDialog<DisputeAnswer>(
    context: context,
    builder: (context) => DisputeAnswerDialog(dispute: dispute),
  );

  @override
  State<DisputeAnswerDialog> createState() => _DisputeAnswerDialogState();
}

class _DisputeAnswerDialogState extends State<DisputeAnswerDialog> {
  final _answer = TextEditingController();
  final _points = TextEditingController(text: '0');

  @override
  void dispose() {
    _answer.dispose();
    _points.dispose();
    super.dispose();
  }

  void _done({required bool accept}) {
    final answer = _answer.text.trim();
    final points = accept ? int.tryParse(_points.text.trim()) : 0;
    if (answer.isEmpty || points == null) {
      return;
    }
    Navigator.of(context)
        .pop(DisputeAnswer(accept: accept, answer: answer, points: points));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final dispute = widget.dispute;
    return AlertDialog(
      title: Text(l10n.leagueDisputeAnswer),
      content: SizedBox(
        width: 380,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(dispute.text),
            TextField(
              key: const ValueKey('dispute-answer'),
              controller: _answer,
              minLines: 2,
              maxLines: 6,
              decoration: InputDecoration(labelText: l10n.leagueAnswer),
            ),
            if (dispute.entryId != null)
              TextField(
                key: const ValueKey('dispute-points'),
                controller: _points,
                keyboardType: const TextInputType.numberWithOptions(
                  signed: true,
                ),
                decoration: InputDecoration(labelText: l10n.leagueAnswerPoints),
              ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => _done(accept: false),
          child: Text(l10n.leagueReject),
        ),
        FilledButton(
          key: const ValueKey('dispute-accept'),
          onPressed: () => _done(accept: true),
          child: Text(l10n.leagueAccept),
        ),
      ],
    );
  }
}
