import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'lti_tool.dart';
import 'tool_url_tile.dart';

class ToolUrlsSection extends StatelessWidget {
  const ToolUrlsSection({required this.tool, super.key});

  final LtiTool tool;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return Column(
      children: [
        ToolUrlTile(label: l10n.lmsLoginUrl, url: tool.loginUrl),
        ToolUrlTile(label: l10n.lmsLaunchUrl, url: tool.launchUrl),
        ToolUrlTile(label: l10n.lmsKeysUrl, url: tool.jwksUrl),
      ],
    );
  }
}
