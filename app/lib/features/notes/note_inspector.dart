import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/pane_heading.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_heading.dart';
import 'note_resolver.dart';
import 'note_tree_row.dart';
import 'notes_providers.dart';

class NoteInspector extends ConsumerWidget {
  const NoteInspector({
    required this.note,
    required this.onOpen,
    required this.onFollow,
    required this.onJump,
    super.key,
  });

  final Note note;
  final ValueChanged<String> onOpen;
  final ValueChanged<String> onFollow;
  final ValueChanged<int> onJump;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final palette = context.palette;
    final notes = ref.watch(notesProvider).value ?? const <Note>[];
    final headings = NoteHeading.parse(note.body);
    final backlinks = NoteResolver.backlinks(note, notes);
    final muted = Theme.of(context).textTheme.bodySmall!
        .copyWith(color: palette.muted);
    Widget hint(String text) => Padding(
      padding: const EdgeInsets.symmetric(
        horizontal: Spacing.l,
        vertical: Spacing.xs,
      ),
      child: Text(text, style: muted),
    );
    return ListView(
      padding: const EdgeInsets.only(bottom: Spacing.l),
      children: [
        PaneHeading(title: l10n.notesOutline),
        if (headings.isEmpty) hint(l10n.notesOutlineNone),
        for (final heading in headings)
          NoteTreeRow(
            depth: heading.level - 1,
            icon: Icons.short_text,
            label: heading.title,
            onTap: () => onJump(heading.offset),
          ),
        PaneHeading(title: l10n.notesLinksIn),
        if (backlinks.isEmpty) hint(l10n.notesLinksNone),
        for (final other in backlinks)
          NoteTreeRow(
            depth: 0,
            icon: Icons.call_received,
            label: other.meta.name,
            onTap: () => onOpen(other.id),
          ),
        PaneHeading(title: l10n.notesLinksOut),
        for (final link in note.meta.links)
          NoteTreeRow(
            depth: 0,
            icon: NoteResolver.resolve(link, notes) == null
                ? Icons.add_link
                : Icons.north_east,
            label: link,
            onTap: () => onFollow(link),
          ),
        if (note.meta.tags.isNotEmpty) ...[
          PaneHeading(title: l10n.notesTags),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: Spacing.l),
            child: Wrap(
              spacing: Spacing.xs,
              runSpacing: Spacing.xs,
              children: [
                for (final tag in note.meta.tags)
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: Spacing.s,
                      vertical: 2,
                    ),
                    decoration: BoxDecoration(
                      color: palette.accentSoft,
                      borderRadius: context.tokens.borderRadiusSmall,
                    ),
                    child: Text(
                      '#$tag',
                      style: muted.copyWith(color: palette.accent),
                    ),
                  ),
              ],
            ),
          ),
        ],
      ],
    );
  }
}
