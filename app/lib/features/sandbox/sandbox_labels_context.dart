import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'sandbox_labels.dart';

extension SandboxLabelsContext on BuildContext {
  SandboxLabels get sandboxLabels =>
      SandboxLabels(l10n, Localizations.localeOf(this).toLanguageTag());
}
