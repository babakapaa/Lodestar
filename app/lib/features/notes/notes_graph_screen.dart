import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../l10n/l10n_context.dart';
import 'notes_graph_view.dart';

class NotesGraphScreen extends StatelessWidget {
  const NotesGraphScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.notesGraph)),
      body: NotesGraphView(onOpen: (id) => context.push(AppRoutes.note(id))),
    );
  }
}
