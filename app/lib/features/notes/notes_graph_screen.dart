import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/empty_state.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'notes_graph.dart';
import 'notes_graph_painter.dart';
import 'notes_providers.dart';

class NotesGraphScreen extends ConsumerWidget {
  const NotesGraphScreen({super.key});

  static const _margin = 120.0;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final notes = ref.watch(notesProvider).value ?? const <Note>[];
    final graph = NotesGraph.of(notes);
    final bounds = graph.bounds.inflate(_margin);
    final origin = -bounds.topLeft;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.notesGraph)),
      body: graph.nodes.isEmpty
          ? EmptyState(
              icon: Icons.hub_outlined,
              title: l10n.emptyNotesTitle,
              body: l10n.emptyNotesBody,
            )
          : InteractiveViewer(
              constrained: false,
              boundaryMargin: const EdgeInsets.all(_margin),
              minScale: 0.2,
              maxScale: 3,
              child: GestureDetector(
                onTapUp: (details) {
                  final node = graph.hit(details.localPosition - origin, 20);
                  if (node != null) {
                    context.push(AppRoutes.note(node.id));
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
