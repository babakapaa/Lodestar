import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'navigator_labels.dart';

extension NavigatorLabelsContext on BuildContext {
  NavigatorLabels get navigatorLabels => NavigatorLabels(l10n);
}
