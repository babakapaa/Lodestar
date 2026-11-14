import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';

class LinkSuggestions extends StatelessWidget {
  const LinkSuggestions({
    required this.notes,
    required this.selected,
    required this.onPick,
    super.key,
  });

  final List<Note> notes;
  final int selected;
  final ValueChanged<Note> onPick;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final textTheme = Theme.of(context).textTheme;
    return Material(
      color: tokens.palette.surface,
      elevation: 4,
      shape: RoundedRectangleBorder(
        borderRadius: tokens.borderRadiusSmall,
        side: BorderSide(color: tokens.palette.border),
      ),
      clipBehavior: Clip.antiAlias,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          for (final (index, note) in notes.indexed)
            InkWell(
              onTap: () => onPick(note),
              child: Container(
                color: index == selected ? tokens.palette.accentSoft : null,
                padding: const EdgeInsets.symmetric(
                  horizontal: Spacing.m,
                  vertical: Spacing.s,
                ),
                child: Row(
                  children: [
                    Icon(
                      Icons.description_outlined,
                      size: 16,
                      color: tokens.palette.muted,
                    ),
                    const SizedBox(width: Spacing.s),
                    Expanded(
                      child: Text(
                        note.meta.name,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    if (note.meta.folder.isNotEmpty)
                      Text(
                        note.meta.folder,
                        style: textTheme.bodySmall!.copyWith(
                          color: tokens.palette.muted,
                        ),
                      ),
                  ],
                ),
              ),
            ),
          const Divider(),
          Padding(
            padding: const EdgeInsets.all(Spacing.s),
            child: Text(
              context.l10n.notesLinkHint,
              style: textTheme.bodySmall!.copyWith(color: tokens.palette.muted),
            ),
          ),
        ],
      ),
    );
  }
}
