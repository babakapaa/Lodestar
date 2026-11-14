import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'hardware_check.dart';
import 'setup_labels.dart';

class HardwareTile extends StatelessWidget {
  const HardwareTile({required this.check, super.key});

  final HardwareCheck check;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = SetupLabels(l10n);
    final palette = context.palette;
    final (icon, color, state) = switch (check.status) {
      HardwareCheck.ok => (
        Icons.check_circle_outline,
        palette.accent,
        l10n.setupHardwareOk,
      ),
      HardwareCheck.failed => (
        Icons.error_outline,
        Theme.of(context).colorScheme.error,
        l10n.setupHardwareFailed,
      ),
      _ => (Icons.remove_circle_outline, palette.muted, l10n.setupHardwareOff),
    };
    return ListTile(
      key: ValueKey('hardware-${check.name}'),
      contentPadding: EdgeInsets.zero,
      leading: Icon(icon, color: color),
      title: Text(labels.hardware(check.name)),
      subtitle: Text(
        [
          state,
          if (check.required) l10n.setupHardwareRequired,
          if (check.detail.isNotEmpty) check.detail,
        ].join('. '),
      ),
    );
  }
}
