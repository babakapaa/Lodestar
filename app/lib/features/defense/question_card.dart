import 'package:material_ui/material_ui.dart';

import '../../design/components/choice_group.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'defense_labels.dart';
import 'defense_question.dart';

class QuestionCard extends StatefulWidget {
  const QuestionCard({
    required this.index,
    required this.question,
    required this.onChanged,
    required this.onRemove,
    super.key,
  });

  final int index;
  final DefenseQuestion question;
  final ValueChanged<DefenseQuestion> onChanged;
  final VoidCallback onRemove;

  @override
  State<QuestionCard> createState() => _QuestionCardState();
}

class _QuestionCardState extends State<QuestionCard> {
  late final _note = TextEditingController(text: widget.question.note);

  @override
  void dispose() {
    _note.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = DefenseLabels(l10n);
    final q = widget.question;
    final textTheme = Theme.of(context).textTheme;
    return LodestarCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      labels.kind(q.kind),
                      style: textTheme.labelMedium!.copyWith(
                        color: context.palette.accent,
                      ),
                    ),
                    Text(
                      '${widget.index + 1}. ${q.text}',
                      style: textTheme.titleSmall,
                    ),
                  ],
                ),
              ),
              IconButton(
                tooltip: l10n.defenseRemoveQuestion,
                icon: const Icon(Icons.close),
                onPressed: widget.onRemove,
              ),
            ],
          ),
          const SizedBox(height: Spacing.s),
          ChoiceGroup<int>(
            values: const [0, 1, 2, 3],
            selected: q.mark ?? -1,
            labelOf: labels.mark,
            onSelected: (mark) => widget.onChanged(q.copyWith(mark: mark)),
          ),
          TextField(
            key: ValueKey('note-${widget.index}'),
            controller: _note,
            minLines: 1,
            maxLines: 4,
            maxLength: 2000,
            decoration: InputDecoration(labelText: l10n.defenseNote),
            onChanged: (note) => widget.onChanged(q.copyWith(note: note)),
          ),
        ],
      ),
    );
  }
}
