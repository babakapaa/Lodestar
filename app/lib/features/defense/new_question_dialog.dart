import 'package:material_ui/material_ui.dart';

import '../../design/components/choice_group.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'defense_labels.dart';
import 'defense_question.dart';
import 'question_kind.dart';

class NewQuestionDialog extends StatefulWidget {
  const NewQuestionDialog({super.key});

  static Future<DefenseQuestion?> show(BuildContext context) =>
      showDialog<DefenseQuestion>(
        context: context,
        builder: (context) => const NewQuestionDialog(),
      );

  @override
  State<NewQuestionDialog> createState() => _NewQuestionDialogState();
}

class _NewQuestionDialogState extends State<NewQuestionDialog> {
  final _text = TextEditingController();
  QuestionKind _kind = QuestionKind.understanding;

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  void _done() {
    final text = _text.text.trim();
    if (text.isNotEmpty) {
      Navigator.of(context).pop(DefenseQuestion(kind: _kind, text: text));
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = DefenseLabels(l10n);
    return AlertDialog(
      title: Text(l10n.defenseAddQuestion),
      content: SizedBox(
        width: 480,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            ChoiceGroup<QuestionKind>(
              values: QuestionKind.values,
              selected: _kind,
              labelOf: labels.kind,
              onSelected: (kind) => setState(() => _kind = kind),
            ),
            const SizedBox(height: Spacing.m),
            TextField(
              key: const ValueKey('question-text'),
              controller: _text,
              autofocus: true,
              minLines: 2,
              maxLines: 5,
              maxLength: 2000,
              decoration: InputDecoration(
                labelText: l10n.defenseQuestionText,
                hintText: labels.example(_kind),
              ),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(onPressed: _done, child: Text(l10n.defenseAdd)),
      ],
    );
  }
}
