import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/layout_size.dart';
import '../router/app_tab.dart';
import 'side_panel.dart';

class DetailFrame extends StatelessWidget {
  const DetailFrame({required this.location, required this.child, super.key});

  final String location;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    final size = LayoutSize.of(MediaQuery.sizeOf(context).width);
    if (size == LayoutSize.compact) {
      return child;
    }
    return Scaffold(
      body: Row(
        children: [
          SidePanel(
            currentIndex: AppTab.of(location)?.index,
            onSelected: (index) => context.go(AppTab.values[index].path),
            extended: size == LayoutSize.expanded,
          ),
          const VerticalDivider(width: 1),
          Expanded(child: child),
        ],
      ),
    );
  }
}
