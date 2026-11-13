import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import '../study/study_labels_context.dart';
import 'team_labels.dart';

extension TeamLabelsContext on BuildContext {
  TeamLabels get teamLabels => TeamLabels(l10n, studyLabels);
}
