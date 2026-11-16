import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import 'side_panel.dart';

class ShellFrame extends StatelessWidget {
  const ShellFrame({
    required this.currentIndex,
    required this.onSelected,
    required this.extended,
    required this.child,
    super.key,
  });

  final int? currentIndex;
  final ValueChanged<int> onSelected;
  final bool extended;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    return Scaffold(
      backgroundColor: palette.frame,
      body: Row(
        children: [
          SidePanel(
            currentIndex: currentIndex,
            onSelected: onSelected,
            extended: extended,
          ),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(
                0,
                Spacing.s,
                Spacing.s,
                Spacing.s,
              ),
              child: DecoratedBox(
                decoration: BoxDecoration(
                  color: palette.bg,
                  borderRadius: tokens.borderRadius,
                  border: Border.all(color: palette.frameEdge),
                  boxShadow: palette.raised,
                ),
                child: ClipRRect(
                  borderRadius: tokens.borderRadius,
                  child: child,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
