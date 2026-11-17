import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_mark.dart';
import '../../design/tokens/tokens_context.dart';

class NavigatorAvatar extends StatelessWidget {
  const NavigatorAvatar({super.key});

  @override
  Widget build(BuildContext context) {
    return LodestarMark(size: 30, color: context.palette.accent);
  }
}
