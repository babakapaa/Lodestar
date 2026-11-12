import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'privacy_providers.dart';

class DeletionSection extends ConsumerStatefulWidget {
  const DeletionSection({super.key});

  @override
  ConsumerState<DeletionSection> createState() => _DeletionSectionState();
}

class _DeletionSectionState extends ConsumerState<DeletionSection> {
  bool _busy = false;
  String? _error;

  Future<bool> _confirm() async {
    final l10n = context.l10n;
    final answer = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.deletionConfirmTitle),
        content: Text(l10n.deletionConfirmBody),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: Text(l10n.actionCancel),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: Text(l10n.deletionRequest),
          ),
        ],
      ),
    );
    return answer ?? false;
  }

  Future<void> _run(Future<void> Function() action) async {
    final l10n = context.l10n;
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await action();
      ref.invalidate(deletionProvider);
    } on ApiException catch (error) {
      _error = ApiErrorText.of(l10n, error);
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  Future<void> _request() async {
    if (await _confirm()) {
      await _run(ref.read(privacyRepositoryProvider).requestDeletion);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final deletion = ref.watch(deletionProvider);
    final error = _error;
    final repository = ref.read(privacyRepositoryProvider);
    return Section(
      title: l10n.deletionTitle,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ...switch (deletion) {
            AsyncData(value: final planned?) => [
              Text(
                l10n.deletionPlanned(
                  DateFormat.yMMMMd(
                    Localizations.localeOf(context).toLanguageTag(),
                  ).format(planned.dueAt.toLocal()),
                ),
              ),
              const SizedBox(height: Spacing.m),
              LodestarButton(
                label: l10n.deletionCancel,
                onPressed: _busy ? null : () => _run(repository.cancelDeletion),
              ),
            ],
            AsyncData() => [
              Text(l10n.deletionBody),
              const SizedBox(height: Spacing.m),
              LodestarButton(
                label: l10n.deletionRequest,
                icon: Icons.delete_outline,
                variant: LodestarButtonVariant.secondary,
                onPressed: _busy ? null : _request,
              ),
            ],
            AsyncError(:final error) => [
              FormError(message: ApiErrorText.of(l10n, error)),
            ],
            _ => [const Center(child: CircularProgressIndicator())],
          },
          if (error != null) FormError(message: error),
        ],
      ),
    );
  }
}
