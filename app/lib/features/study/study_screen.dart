import 'package:material_ui/material_ui.dart';

import '../../core/navigation/tab_page.dart';
import '../../design/components/empty_state.dart';
import '../../l10n/l10n_context.dart';

class StudyScreen extends StatelessWidget {
  const StudyScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return TabPage(
      title: l10n.tabStudy,
      child: EmptyState(
        icon: Icons.school_outlined,
        title: l10n.emptyStudyTitle,
        body: l10n.emptyStudyBody,
      ),
    );
  }
}
