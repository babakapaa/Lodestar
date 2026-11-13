import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import '../study/course_person.dart';

abstract final class TeamPersonPicker {
  static Future<CoursePerson?> show(
    BuildContext context, {
    required String title,
    required List<CoursePerson> people,
  }) => showDialog<CoursePerson>(
    context: context,
    builder: (context) => SimpleDialog(
      title: Text(title),
      children: [
        if (people.isEmpty)
          Padding(
            padding: const EdgeInsets.all(24),
            child: Text(context.l10n.teamPickNobody),
          ),
        for (final person in people)
          SimpleDialogOption(
            key: ValueKey('pick-${person.userId}'),
            onPressed: () => Navigator.of(context).pop(person),
            child: Text(person.label),
          ),
      ],
    ),
  );
}
