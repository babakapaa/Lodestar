import 'package:flutter/services.dart';
import 'package:material_ui/material_ui.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'run_kind.dart';
import 'sandbox_run.dart';

class RunLink extends StatelessWidget {
  const RunLink({required this.run, super.key});

  final SandboxRun run;

  Future<void> _copy(BuildContext context) async {
    final messenger = ScaffoldMessenger.of(context);
    final l10n = context.l10n;
    await Clipboard.setData(ClipboardData(text: run.url));
    messenger.showSnackBar(SnackBar(content: Text(l10n.sandboxCopied)));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Wrap(
          spacing: Spacing.s,
          children: [
            FilledButton.icon(
              key: const ValueKey('run-open'),
              onPressed: () => launchUrl(Uri.parse(run.url)),
              icon: const Icon(Icons.open_in_new),
              label: Text(l10n.sandboxOpen),
            ),
            OutlinedButton.icon(
              key: const ValueKey('run-copy'),
              onPressed: () => _copy(context),
              icon: const Icon(Icons.copy_outlined),
              label: Text(l10n.sandboxCopy),
            ),
          ],
        ),
        if (run.kind == RunKind.demo) ...[
          const SizedBox(height: Spacing.l),
          Text(l10n.sandboxDemoQr),
          const SizedBox(height: Spacing.s),
          ColoredBox(
            color: Colors.white,
            child: QrImageView(
              key: const ValueKey('run-qr'),
              data: run.url,
              size: 200,
            ),
          ),
        ],
      ],
    );
  }
}
