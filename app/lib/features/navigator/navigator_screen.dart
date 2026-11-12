import 'package:material_ui/material_ui.dart';

import '../../core/navigation/tab_page.dart';
import '../../design/components/empty_state.dart';
import '../../l10n/l10n_context.dart';

class NavigatorScreen extends StatelessWidget {
  const NavigatorScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return TabPage(
      title: l10n.tabNavigator,
      child: EmptyState(
        icon: Icons.explore_outlined,
        title: l10n.emptyNavigatorTitle,
        body: l10n.emptyNavigatorBody,
      ),
    );
  }
}
