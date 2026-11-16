import 'package:material_ui/material_ui.dart';

import '../../design/components/brand_mark.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../router/app_tab.dart';
import 'side_account.dart';
import 'side_nav_item.dart';

class SideNav extends StatelessWidget {
  const SideNav({
    required this.currentIndex,
    required this.onSelected,
    super.key,
  });

  static const width = 248.0;

  final int? currentIndex;
  final ValueChanged<int> onSelected;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return ColoredBox(
      color: context.tokens.palette.frame,
      child: SizedBox(
        width: width,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(
            Spacing.m,
            Spacing.l,
            Spacing.m,
            Spacing.m,
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Padding(
                padding: EdgeInsets.symmetric(horizontal: Spacing.s),
                child: BrandMark(showName: true),
              ),
              const SizedBox(height: Spacing.xl),
              for (final tab in AppTab.values)
                SideNavItem(
                  icon: tab.icon,
                  selectedIcon: tab.selectedIcon,
                  label: tab.label(l10n),
                  selected: currentIndex == tab.index,
                  onTap: () => onSelected(tab.index),
                ),
              const Spacer(),
              const SideAccount(),
            ],
          ),
        ),
      ),
    );
  }
}
