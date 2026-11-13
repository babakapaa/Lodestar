import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import '../study/study_labels_context.dart';
import 'league_labels.dart';

extension LeagueLabelsContext on BuildContext {
  LeagueLabels get leagueLabels => LeagueLabels(l10n, studyLabels);
}
