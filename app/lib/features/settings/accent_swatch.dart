import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';

class AccentSwatch extends StatelessWidget {
  const AccentSwatch({
    required this.color,
    required this.label,
    required this.auto,
    required this.selected,
    required this.onTap,
    super.key,
  });

  final Color color;
  final String label;
  final bool auto;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    return Tooltip(
      message: label,
      child: Semantics(
        button: true,
        selected: selected,
        label: label,
        child: InkResponse(
          onTap: onTap,
          radius: 22,
          child: Container(
            width: 36,
            height: 36,
            padding: const EdgeInsets.all(3),
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              border: Border.all(
                color: selected ? palette.text : palette.border,
                width: 2,
              ),
            ),
            child: DecoratedBox(
              decoration: BoxDecoration(color: color, shape: BoxShape.circle),
              child: auto
                  ? Icon(Icons.auto_awesome, size: 14, color: palette.surface)
                  : null,
            ),
          ),
        ),
      ),
    );
  }
}
