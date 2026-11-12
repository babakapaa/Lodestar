import 'package:material_ui/material_ui.dart';

import 'lodestar_tokens.dart';
import 'palette.dart';

extension TokensContext on BuildContext {
  LodestarTokens get tokens => Theme.of(this).extension<LodestarTokens>()!;

  Palette get palette => tokens.palette;
}
