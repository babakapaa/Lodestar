import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';
import 'role_color.dart';

class RoleDot extends StatelessWidget {
  const RoleDot({required this.color, this.size = 12, super.key});

  final String color;
  final double size;

  @override
  Widget build(BuildContext context) {
    final parsed = RoleColor.parse(color);
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: parsed ?? Colors.transparent,
        border: Border.all(color: parsed ?? context.tokens.palette.muted),
      ),
    );
  }
}
