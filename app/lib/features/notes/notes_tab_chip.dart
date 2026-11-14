import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class NotesTabChip extends StatelessWidget {
  const NotesTabChip({
    required this.icon,
    required this.label,
    required this.active,
    required this.onTap,
    required this.onClose,
    super.key,
  });

  final IconData icon;
  final String label;
  final bool active;
  final VoidCallback onTap;
  final VoidCallback onClose;

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    return Material(
      color: active ? palette.surface : Colors.transparent,
      child: InkWell(
        onTap: onTap,
        child: Container(
          constraints: const BoxConstraints(minWidth: 120, maxWidth: 220),
          padding: const EdgeInsets.only(left: Spacing.m, right: Spacing.xs),
          decoration: BoxDecoration(
            border: Border(
              top: BorderSide(
                color: active ? palette.accent : Colors.transparent,
                width: 2,
              ),
              right: BorderSide(color: palette.border),
            ),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                icon,
                size: 16,
                color: active ? palette.accent : palette.muted,
              ),
              const SizedBox(width: Spacing.s),
              Flexible(
                child: Text(
                  label,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: Theme.of(context).textTheme.bodyMedium!
                      .copyWith(color: active ? palette.text : palette.muted),
                ),
              ),
              IconButton(
                tooltip: context.l10n.notesCloseTab,
                visualDensity: VisualDensity.compact,
                iconSize: 16,
                icon: const Icon(Icons.close),
                onPressed: onClose,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
