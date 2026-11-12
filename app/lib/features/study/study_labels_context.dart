import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'study_labels.dart';

extension StudyLabelsContext on BuildContext {
  StudyLabels get studyLabels =>
      StudyLabels(l10n, Localizations.localeOf(this).toLanguageTag());
}
