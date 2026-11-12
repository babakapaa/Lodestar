import 'package:material_ui/material_ui.dart';

import '../../core/navigation/tab_page.dart';
import '../../design/components/empty_state.dart';
import '../../l10n/l10n_context.dart';

class NotesScreen extends StatelessWidget {
  const NotesScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return TabPage(
      title: l10n.tabNotes,
      child: EmptyState(
        icon: Icons.sticky_note_2_outlined,
        title: l10n.emptyNotesTitle,
        body: l10n.emptyNotesBody,
      ),
    );
  }
}
