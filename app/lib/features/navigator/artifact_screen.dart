import 'package:material_ui/material_ui.dart';

import 'artifact_panel.dart';
import 'navigator_artifact.dart';

class ArtifactScreen extends StatelessWidget {
  const ArtifactScreen({required this.artifact, super.key});

  final NavigatorArtifact artifact;

  static Future<void> show(BuildContext context, NavigatorArtifact artifact) =>
      Navigator.of(context).push(
        MaterialPageRoute<void>(
          fullscreenDialog: true,
          builder: (context) => ArtifactScreen(artifact: artifact),
        ),
      );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: ArtifactPanel(
          artifact: artifact,
          onClose: () => Navigator.of(context).pop(),
        ),
      ),
    );
  }
}
