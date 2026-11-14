import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';

class NavigatorAvatar extends StatelessWidget {
  const NavigatorAvatar({super.key});

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    return Container(
      width: 30,
      height: 30,
      decoration: BoxDecoration(
        color: palette.accent,
        borderRadius: context.tokens.borderRadiusSmall,
      ),
      child: Icon(Icons.explore, size: 18, color: palette.onAccent),
    );
  }
}
