import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/empty_state.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'notes_graph.dart';
import 'notes_graph_fit.dart';
import 'notes_graph_painter.dart';
import 'notes_providers.dart';
import 'note_actions.dart';

class NotesGraphView extends ConsumerWidget {
  const NotesGraphView({required this.onOpen, super.key});

  static const _margin = 80.0;

  final ValueChanged<String> onOpen;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final notes = ref.watch(notesProvider).value ?? const <Note>[];
    final graph = NotesGraph.of(notes);
    final bounds = graph.bounds.inflate(_margin);
    final origin = -bounds.topLeft;
    if (graph.nodes.isEmpty) {
      return EmptyState(
        icon: Icons.hub_outlined,
        title: l10n.emptyNotesTitle,
        body: l10n.emptyNotesBody,
        actionLabel: l10n.notesNew,
        onAction: () => NoteActions.create(context, ref),
      );
    }
    return LayoutBuilder(
      builder: (context, constraints) => NotesGraphFit(
        viewport: constraints.biggest,
        content: bounds.size,
        child: GestureDetector(
          onTapUp: (details) {
            final node = graph.hit(details.localPosition - origin, 20);
            if (node != null) {
              onOpen(node.id);
            }
          },
          child: Semantics(
            label: l10n.notesGraphLabel(graph.nodes.length),
            child: CustomPaint(
              size: bounds.size,
              painter: NotesGraphPainter(
                graph: graph,
                palette: context.palette,
                labelStyle: Theme.of(context).textTheme.labelSmall!,
                origin: origin,
              ),
            ),
          ),
        ),
      ),
    );
  }
}
