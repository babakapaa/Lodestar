import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'defense_labels.dart';
import 'question_kind.dart';

class QuestionTypesHint extends StatelessWidget {
  const QuestionTypesHint({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = DefenseLabels(l10n);
    return ExpansionTile(
      tilePadding: EdgeInsets.zero,
      title: Text(l10n.defenseQuestionTypes),
      children: [
        for (final kind in QuestionKind.values)
          ListTile(
            contentPadding: EdgeInsets.zero,
            title: Text(labels.kind(kind)),
            subtitle: Text(labels.example(kind)),
          ),
      ],
    );
  }
}
