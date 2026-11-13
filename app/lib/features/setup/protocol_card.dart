import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/files/files_providers.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'setup_overview.dart';
import 'setup_providers.dart';

class ProtocolCard extends ConsumerWidget {
  const ProtocolCard({required this.overview, super.key});

  final SetupOverview overview;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final api = ref.read(setupApiProvider);
    return Section(
      title: l10n.setupProtocol,
      child: ListTile(
        contentPadding: EdgeInsets.zero,
        leading: const Icon(Icons.description_outlined),
        title: Text(
          overview.signed
              ? l10n.setupProtocolSigned(overview.protocolVersion)
              : l10n.setupProtocolUnsigned,
        ),
        trailing: Wrap(
          children: [
            for (final format in const ['pdf', 'md'])
              IconButton(
                key: ValueKey('setup-protocol-$format'),
                tooltip: format == 'pdf'
                    ? l10n.setupDownloadPdf
                    : l10n.setupDownload,
                icon: Icon(
                  format == 'pdf'
                      ? Icons.picture_as_pdf_outlined
                      : Icons.download_outlined,
                ),
                onPressed: () => TeamAction.run(context, () async {
                  final bytes = await api.protocol(format);
                  await ref
                      .read(exportSaverProvider)
                      .save('protocol', bytes, extension: format);
                }),
              ),
            if (!overview.signed)
              TextButton(
                key: const ValueKey('setup-sign'),
                onPressed: () async {
                  await TeamAction.run(context, api.sign);
                  ref.invalidate(setupOverviewProvider);
                },
                child: Text(l10n.setupSign),
              ),
          ],
        ),
      ),
    );
  }
}
