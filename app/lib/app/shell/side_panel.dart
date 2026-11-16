import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../router/app_tab.dart';
import '../../design/components/brand_mark.dart';
import '../../core/navigation/profile_button.dart';
import 'side_nav.dart';

class SidePanel extends StatelessWidget {
  const SidePanel({
    required this.currentIndex,
    required this.onSelected,
    required this.extended,
    super.key,
  });

  final int? currentIndex;
  final ValueChanged<int> onSelected;
  final bool extended;

  @override
  Widget build(BuildContext context) {
    if (extended) {
      return SideNav(currentIndex: currentIndex, onSelected: onSelected);
    }
    final l10n = context.l10n;
    return NavigationRail(
      extended: extended,
      minExtendedWidth: 248,
      selectedIndex: currentIndex,
      onDestinationSelected: onSelected,
      labelType: extended
          ? NavigationRailLabelType.none
          : NavigationRailLabelType.all,
      leading: Padding(
        padding: const EdgeInsets.symmetric(vertical: Spacing.l),
        child: BrandMark(showName: extended, width: 216),
      ),
      trailing: Expanded(
        child: Align(
          alignment: Alignment.bottomCenter,
          child: Padding(
            padding: const EdgeInsets.only(bottom: Spacing.l),
            child: ProfileButton(showLabel: extended),
          ),
        ),
      ),
      destinations: [
        for (final tab in AppTab.values)
          NavigationRailDestination(
            icon: Icon(tab.icon),
            selectedIcon: Icon(tab.selectedIcon),
            label: Text(tab.label(l10n)),
          ),
      ],
    );
  }
}
