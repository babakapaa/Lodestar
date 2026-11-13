import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'setup_labels.dart';
import 'setup_legal_answer.dart';
import 'setup_providers.dart';

class LegalQuestionCard extends ConsumerStatefulWidget {
  const LegalQuestionCard({required this.answer, super.key});

  final SetupLegalAnswer answer;

  @override
  ConsumerState<LegalQuestionCard> createState() => _LegalQuestionCardState();
}

class _LegalQuestionCardState extends ConsumerState<LegalQuestionCard> {
  late final _answer = TextEditingController(text: widget.answer.answer);
  late final _basis = TextEditingController(text: widget.answer.basis);

  @override
  void dispose() {
    _answer.dispose();
    _basis.dispose();
    super.dispose();
  }

  Future<void> _run(Future<void> Function() action) async {
    await TeamAction.run(context, action);
    ref.invalidate(setupOverviewProvider);
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final a = widget.answer;
    final api = ref.read(setupApiProvider);
    final state = a.approved
        ? l10n.setupAnswerApproved
        : a.answered
        ? l10n.setupAnswerWaiting
        : l10n.setupAnswerEmpty;
    return Card(
      margin: const EdgeInsets.only(bottom: Spacing.m),
      child: Padding(
        padding: const EdgeInsets.all(Spacing.l),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              SetupLabels(l10n).question(a.question, a.title),
              style: Theme.of(context).textTheme.titleSmall,
            ),
            const SizedBox(height: Spacing.xs),
            Text(state),
            const SizedBox(height: Spacing.m),
            TextField(
              key: ValueKey('legal-answer-${a.question}'),
              controller: _answer,
              maxLines: null,
              decoration: InputDecoration(labelText: l10n.setupAnswer),
            ),
            const SizedBox(height: Spacing.s),
            TextField(
              key: ValueKey('legal-basis-${a.question}'),
              controller: _basis,
              decoration: InputDecoration(labelText: l10n.setupBasis),
            ),
            const SizedBox(height: Spacing.m),
            Wrap(
              spacing: Spacing.m,
              children: [
                TextButton(
                  key: ValueKey('legal-save-${a.question}'),
                  onPressed: () => _run(
                    () => api.answer(a.question, _answer.text, _basis.text),
                  ),
                  child: Text(l10n.actionSave),
                ),
                if (a.answered && !a.approved)
                  TextButton(
                    key: ValueKey('legal-approve-${a.question}'),
                    onPressed: () => _run(() => api.approve(a.question)),
                    child: Text(l10n.setupApprove),
                  ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
