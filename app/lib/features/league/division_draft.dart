import 'package:material_ui/material_ui.dart';

import 'division.dart';

class DivisionDraft {
  DivisionDraft(Division division)
    : name = TextEditingController(text: division.name),
      from = TextEditingController(text: '${division.from}'),
      gpuHours = TextEditingController(text: '${division.gpuHours}');

  final TextEditingController name;
  final TextEditingController from;
  final TextEditingController gpuHours;

  Division? get value {
    final start = int.tryParse(from.text.trim());
    final hours = int.tryParse(gpuHours.text.trim());
    final title = name.text.trim();
    if (start == null || hours == null || title.isEmpty) {
      return null;
    }
    return Division(name: title, from: start, gpuHours: hours);
  }

  void dispose() {
    name.dispose();
    from.dispose();
    gpuHours.dispose();
  }
}
