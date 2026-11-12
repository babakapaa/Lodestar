import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../router/app_tab.dart';

class BottomNav extends StatelessWidget {
  const BottomNav({
    required this.currentIndex,
    required this.onSelected,
    super.key,
  });

  final int currentIndex;
  final ValueChanged<int> onSelected;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return DecoratedBox(
      decoration: BoxDecoration(
        border: Border(top: BorderSide(color: context.palette.border)),
      ),
      child: NavigationBar(
        selectedIndex: currentIndex,
        onDestinationSelected: onSelected,
        destinations: [
          for (final tab in AppTab.values)
            NavigationDestination(
              icon: Icon(tab.icon),
              selectedIcon: Icon(tab.selectedIcon),
              label: tab.label(l10n),
            ),
        ],
      ),
    );
  }
}
