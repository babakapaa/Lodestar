import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'league_solution.dart';
import 'solution_review.dart';
import 'task_details.dart';

class ReviewDialog extends StatefulWidget {
  const ReviewDialog({
    required this.details,
    required this.solution,
    super.key,
  });

  final TaskDetails details;
  final LeagueSolution solution;

  static Future<SolutionReview?> show(
    BuildContext context,
    TaskDetails details,
    LeagueSolution solution,
  ) => showDialog<SolutionReview>(
    context: context,
    builder: (context) => ReviewDialog(details: details, solution: solution),
  );

  @override
  State<ReviewDialog> createState() => _ReviewDialogState();
}

class _ReviewDialogState extends State<ReviewDialog> {
  late final _points = TextEditingController(
    text: '${widget.solution.points ?? widget.details.task.points}',
  );
  late final _comment = TextEditingController(text: widget.solution.comment);
  late final Set<String> _passed = {
    ...(widget.solution.reviewed
        ? widget.solution.passed
        : widget.solution.members),
  };

  @override
  void dispose() {
    _points.dispose();
    _comment.dispose();
    super.dispose();
  }

  void _done() {
    final points = int.tryParse(_points.text.trim());
    if (points == null || points < 0) {
      return;
    }
    Navigator.of(context).pop(
      SolutionReview(
        points: points,
        passed: [
          for (final id in widget.solution.members)
            if (_passed.contains(id)) id,
        ],
        comment: _comment.text.trim(),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(widget.solution.teamName),
      content: SizedBox(
        width: 380,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              TextField(
                key: const ValueKey('review-points'),
                controller: _points,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(labelText: l10n.leagueReviewPoints),
              ),
              Text(l10n.leagueReviewPassed),
              for (final id in widget.solution.members)
                CheckboxListTile(
                  key: ValueKey('passed-$id'),
                  contentPadding: EdgeInsets.zero,
                  title: Text(widget.details.nameOf(id)),
                  value: _passed.contains(id),
                  onChanged: (value) => setState(
                    () => value ?? false ? _passed.add(id) : _passed.remove(id),
                  ),
                ),
              TextField(
                controller: _comment,
                minLines: 2,
                maxLines: 6,
                decoration: InputDecoration(
                  labelText: l10n.leagueReviewComment,
                ),
              ),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(
          key: const ValueKey('review-save'),
          onPressed: _done,
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
