import 'package:material_ui/material_ui.dart';

import '../../../design/tokens/spacing.dart';
import '../../../design/tokens/tokens_context.dart';
import '../../../l10n/l10n_context.dart';

class CanvasToolbar extends StatelessWidget {
  const CanvasToolbar({
    required this.readOnly,
    required this.onAddCard,
    required this.onAddNote,
    required this.onZoomIn,
    required this.onZoomOut,
    required this.onFit,
    super.key,
  });

  final bool readOnly;
  final VoidCallback onAddCard;
  final VoidCallback onAddNote;
  final VoidCallback onZoomIn;
  final VoidCallback onZoomOut;
  final VoidCallback onFit;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final tokens = context.tokens;
    return Material(
      color: tokens.palette.surface,
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: tokens.borderRadius,
        side: BorderSide(color: tokens.palette.border),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: Spacing.xs),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (!readOnly) ...[
              IconButton(
                tooltip: l10n.canvasAddCard,
                icon: const Icon(Icons.note_add_outlined),
                onPressed: onAddCard,
              ),
              IconButton(
                tooltip: l10n.canvasAddNote,
                icon: const Icon(Icons.description_outlined),
                onPressed: onAddNote,
              ),
              const SizedBox(height: 24, child: VerticalDivider()),
            ],
            IconButton(
              tooltip: l10n.canvasZoomOut,
              icon: const Icon(Icons.zoom_out),
              onPressed: onZoomOut,
            ),
            IconButton(
              tooltip: l10n.canvasZoomIn,
              icon: const Icon(Icons.zoom_in),
              onPressed: onZoomIn,
            ),
            IconButton(
              tooltip: l10n.canvasFit,
              icon: const Icon(Icons.fit_screen_outlined),
              onPressed: onFit,
            ),
          ],
        ),
      ),
    );
  }
}
