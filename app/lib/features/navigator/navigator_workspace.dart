import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';
import 'artifact_panel.dart';
import 'dialog_screen.dart';
import 'navigator_artifact.dart';
import 'navigator_history.dart';

class NavigatorWorkspace extends StatefulWidget {
  const NavigatorWorkspace({super.key});

  static const historyWidth = 280.0;
  static const artifactWidth = 460.0;

  @override
  State<NavigatorWorkspace> createState() => _NavigatorWorkspaceState();
}

class _NavigatorWorkspaceState extends State<NavigatorWorkspace> {
  String? _opened;
  String? _current;
  int _fresh = 0;
  NavigatorArtifact? _artifact;

  void _open(String id) => setState(() {
    _opened = id;
    _current = id;
    _artifact = null;
  });

  void _new() => setState(() {
    _opened = null;
    _current = null;
    _fresh++;
    _artifact = null;
  });

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    final opened = _opened;
    final artifact = _artifact;
    return Row(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Container(
          width: NavigatorWorkspace.historyWidth,
          color: palette.side,
          child: NavigatorHistory(
            selected: _current,
            onOpen: _open,
            onNew: _new,
          ),
        ),
        const VerticalDivider(width: 1),
        Expanded(
          child: DialogScreen(
            key: ValueKey(opened ?? 'new-$_fresh'),
            dialogId: opened,
            onArtifact: (value) => setState(() => _artifact = value),
            onCreated: (id) => setState(() => _current = id),
            onDeleted: _new,
          ),
        ),
        if (artifact != null) ...[
          const VerticalDivider(width: 1),
          SizedBox(
            width: NavigatorWorkspace.artifactWidth,
            child: ArtifactPanel(
              key: ValueKey(artifact),
              artifact: artifact,
              onClose: () => setState(() => _artifact = null),
            ),
          ),
        ],
      ],
    );
  }
}
