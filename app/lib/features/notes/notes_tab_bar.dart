import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'notes_providers.dart';
import 'notes_tab_chip.dart';
import 'notes_tabs.dart';

class NotesTabBar extends ConsumerWidget {
  const NotesTabBar({required this.trailing, super.key});

  final List<Widget> trailing;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final tabs = ref.watch(notesTabsProvider);
    final controller = ref.read(notesTabsProvider.notifier);
    final notes = ref.watch(notesProvider).value ?? const <Note>[];
    String nameOf(String id) {
      if (id == NotesTabs.graph) {
        return context.l10n.notesGraph;
      }
      for (final note in notes) {
        if (note.id == id) {
          return note.meta.label;
        }
      }
      return context.l10n.notesNotFound;
    }

    return Container(
      height: 40,
      color: context.palette.side,
      child: Row(
        children: [
          Expanded(
            child: ListView(
              scrollDirection: Axis.horizontal,
              children: [
                for (final id in tabs.open)
                  NotesTabChip(
                    key: ValueKey('tab:$id'),
                    icon: id == NotesTabs.graph
                        ? Icons.hub_outlined
                        : notes.any((n) => n.id == id && n.meta.isCanvas)
                        ? Icons.dashboard_outlined
                        : Icons.description_outlined,
                    label: nameOf(id),
                    active: id == tabs.active,
                    onTap: () => controller.activate(id),
                    onClose: () => controller.close(id),
                  ),
              ],
            ),
          ),
          ...trailing,
        ],
      ),
    );
  }
}
