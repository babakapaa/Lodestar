import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/pane_action.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../notes/markdown_controller.dart';
import '../notes/note_actions.dart';
import 'navigator_artifact.dart';

class ArtifactPanel extends ConsumerStatefulWidget {
  const ArtifactPanel({
    required this.artifact,
    required this.onClose,
    super.key,
  });

  final NavigatorArtifact artifact;
  final VoidCallback onClose;

  @override
  ConsumerState<ArtifactPanel> createState() => _ArtifactPanelState();
}

class _ArtifactPanelState extends ConsumerState<ArtifactPanel> {
  late final _text = MarkdownController()..text = widget.artifact.text;

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  Future<void> _copy() async {
    final messenger = ScaffoldMessenger.of(context);
    final copied = context.l10n.letterCopied;
    await Clipboard.setData(ClipboardData(text: _text.text));
    messenger.showSnackBar(SnackBar(content: Text(copied)));
  }

  Future<void> _toNote() async {
    final messenger = ScaffoldMessenger.of(context);
    final saved = context.l10n.artifactSaved;
    await NoteActions.addText(ref, widget.artifact.title, _text.text);
    messenger.showSnackBar(SnackBar(content: Text(saved)));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = context.palette;
    final textTheme = Theme.of(context).textTheme;
    return ColoredBox(
      color: palette.surface,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          SizedBox(
            height: 56,
            child: Row(
              children: [
                const SizedBox(width: Spacing.l),
                Icon(Icons.description_outlined, color: palette.accent),
                const SizedBox(width: Spacing.s),
                Expanded(
                  child: Text(
                    widget.artifact.title,
                    style: textTheme.titleMedium,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                PaneAction(
                  icon: Icons.copy_outlined,
                  tooltip: l10n.artifactCopy,
                  onPressed: _copy,
                ),
                PaneAction(
                  icon: Icons.note_add_outlined,
                  tooltip: l10n.artifactToNote,
                  onPressed: _toNote,
                ),
                PaneAction(
                  icon: Icons.close,
                  tooltip: l10n.artifactClose,
                  onPressed: widget.onClose,
                ),
                const SizedBox(width: Spacing.s),
              ],
            ),
          ),
          const Divider(),
          Padding(
            padding: const EdgeInsets.fromLTRB(
              Spacing.l,
              Spacing.m,
              Spacing.l,
              0,
            ),
            child: Text(
              l10n.artifactHint,
              style: textTheme.bodySmall!.copyWith(color: palette.muted),
            ),
          ),
          Expanded(
            child: TextField(
              controller: _text,
              maxLines: null,
              expands: true,
              textAlignVertical: TextAlignVertical.top,
              style: textTheme.bodyLarge!.copyWith(height: 1.6),
              decoration: const InputDecoration(
                border: InputBorder.none,
                enabledBorder: InputBorder.none,
                focusedBorder: InputBorder.none,
                filled: false,
                contentPadding: EdgeInsets.all(Spacing.l),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
