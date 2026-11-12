import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'people_search.dart';
import 'person_hit.dart';

class PersonPickerDialog extends StatelessWidget {
  const PersonPickerDialog({super.key});

  static Future<PersonHit?> show(BuildContext context) => showDialog<PersonHit>(
    context: context,
    builder: (context) => const PersonPickerDialog(),
  );

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text(context.l10n.peoplePick),
      content: SizedBox(
        width: 420,
        child: SingleChildScrollView(
          child: PeopleSearch(
            onSelected: (person) => Navigator.of(context).pop(person),
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(context.l10n.actionCancel),
        ),
      ],
    );
  }
}
