import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import 'bottom_nav.dart';
import '../../core/navigation/layout_size.dart';
import '../../features/tour/tour_gate.dart';
import 'side_panel.dart';

class AdaptiveShell extends StatelessWidget {
  const AdaptiveShell({required this.shell, super.key});

  final StatefulNavigationShell shell;

  void _select(int index) {
    shell.goBranch(index, initialLocation: index == shell.currentIndex);
  }

  @override
  Widget build(BuildContext context) {
    final size = LayoutSize.of(MediaQuery.sizeOf(context).width);
    if (size == LayoutSize.compact) {
      return TourGate(
        child: Scaffold(
          body: shell,
          bottomNavigationBar: BottomNav(
            currentIndex: shell.currentIndex,
            onSelected: _select,
          ),
        ),
      );
    }
    return TourGate(
      child: Scaffold(
        body: Row(
          children: [
            SidePanel(
              currentIndex: shell.currentIndex,
              onSelected: _select,
              extended: size == LayoutSize.expanded,
            ),
            const VerticalDivider(width: 1),
            Expanded(child: shell),
          ],
        ),
      ),
    );
  }
}
