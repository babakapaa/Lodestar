import 'package:material_ui/material_ui.dart';

import '../../core/navigation/tab_page.dart';
import '../../design/components/empty_state.dart';
import '../../l10n/l10n_context.dart';

class TodayScreen extends StatelessWidget {
  const TodayScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return TabPage(
      title: l10n.tabToday,
      child: EmptyState(
        icon: Icons.wb_sunny_outlined,
        title: l10n.emptyTodayTitle,
        body: l10n.emptyTodayBody,
      ),
    );
  }
}
