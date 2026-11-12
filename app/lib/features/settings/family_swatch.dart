import 'package:material_ui/material_ui.dart';

import '../../design/tokens/theme_family.dart';

class FamilySwatch extends StatelessWidget {
  const FamilySwatch({required this.family, super.key});

  final ThemeFamily family;

  @override
  Widget build(BuildContext context) {
    final palette = family.spec.palette(Theme.of(context).brightness);
    return Container(
      width: 16,
      height: 16,
      decoration: BoxDecoration(
        color: palette.accent,
        shape: BoxShape.circle,
        border: Border.all(color: palette.border, width: 2),
      ),
    );
  }
}
