import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../../core/files/files_providers.dart';
import 'privacy_providers.dart';

class ExportSection extends ConsumerStatefulWidget {
  const ExportSection({super.key});

  @override
  ConsumerState<ExportSection> createState() => _ExportSectionState();
}

class _ExportSectionState extends ConsumerState<ExportSection> {
  bool _busy = false;
  bool _done = false;
  String? _error;

  Future<void> _export() async {
    final l10n = context.l10n;
    setState(() {
      _busy = true;
      _done = false;
      _error = null;
    });
    try {
      final bytes = await ref.read(privacyRepositoryProvider).export();
      final day = DateTime.now().toIso8601String().substring(0, 10);
      await ref.read(exportSaverProvider).save('lodestar-export-$day', bytes);
      _done = true;
    } on Exception catch (error) {
      _error = ApiErrorText.of(l10n, error);
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    return Section(
      title: l10n.exportTitle,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(l10n.exportBody),
          const SizedBox(height: Spacing.m),
          LodestarButton(
            label: _busy ? l10n.exportBusy : l10n.exportAction,
            icon: Icons.download_outlined,
            variant: LodestarButtonVariant.secondary,
            onPressed: _busy ? null : _export,
          ),
          if (_done)
            Padding(
              padding: const EdgeInsets.only(top: Spacing.m),
              child: Text(l10n.exportDone),
            ),
          if (error != null) FormError(message: error),
        ],
      ),
    );
  }
}
