import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/pane_heading.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../notes/note_tree_row.dart';
import 'navigator_dialog.dart';
import 'navigator_providers.dart';

class NavigatorHistory extends ConsumerWidget {
  const NavigatorHistory({
    required this.selected,
    required this.onOpen,
    required this.onNew,
    super.key,
  });

  final String? selected;
  final ValueChanged<String> onOpen;
  final VoidCallback onNew;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final dialogs = ref.watch(dialogsProvider).value?.value ?? const [];
    final now = DateTime.now();
    final yesterday = now.subtract(const Duration(days: 1));
    final groups = <String, List<NavigatorDialog>>{};
    for (final dialog in dialogs) {
      final at = dialog.changedAt.toLocal();
      final label = DateUtils.isSameDay(at, now)
          ? l10n.chatsToday
          : DateUtils.isSameDay(at, yesterday)
          ? l10n.chatsYesterday
          : l10n.navigatorEarlier;
      groups.putIfAbsent(label, () => []).add(dialog);
    }
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(
            Spacing.m,
            Spacing.l,
            Spacing.m,
            Spacing.s,
          ),
          child: FilledButton.icon(
            onPressed: onNew,
            icon: const Icon(Icons.add_comment_outlined),
            label: Text(l10n.navigatorNewQuestion),
          ),
        ),
        Expanded(
          child: ListView(
            children: [
              for (final MapEntry(key: label, value: items)
                  in groups.entries) ...[
                PaneHeading(title: label),
                for (final dialog in items)
                  NoteTreeRow(
                    key: ValueKey('dialog:${dialog.id}'),
                    depth: 0,
                    icon: dialog.scopeKind.isEmpty
                        ? Icons.chat_bubble_outline
                        : Icons.assignment_outlined,
                    label: dialog.title.isEmpty
                        ? l10n.navigatorUntitled
                        : dialog.title,
                    selected: dialog.id == selected,
                    onTap: () => onOpen(dialog.id),
                  ),
              ],
            ],
          ),
        ),
      ],
    );
  }
}
