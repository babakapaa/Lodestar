import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'admin_providers.dart';
import 'audit_tile.dart';

class AuditScreen extends ConsumerWidget {
  const AuditScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final page = ref.watch(auditProvider);
    return Scaffold(
      appBar: AppBar(title: Text(l10n.auditTitle)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: switch (page) {
              AsyncData(:final value) when value.entries.isEmpty => EmptyState(
                icon: Icons.history,
                title: l10n.auditEmptyTitle,
                body: l10n.auditEmptyBody,
              ),
              AsyncData(:final value) => ListView(
                children: [
                  for (final entry in value.entries) AuditTile(entry: entry),
                  if (value.hasMore)
                    Padding(
                      padding: const EdgeInsets.all(Spacing.l),
                      child: LodestarButton(
                        label: l10n.notificationsMore,
                        variant: LodestarButtonVariant.secondary,
                        onPressed: () =>
                            ref.read(auditProvider.notifier).loadMore(),
                      ),
                    ),
                ],
              ),
              AsyncError(:final error) => Padding(
                padding: const EdgeInsets.all(Spacing.l),
                child: FormError(message: ApiErrorText.of(l10n, error)),
              ),
              _ => const Center(child: CircularProgressIndicator()),
            },
          ),
        ),
      ),
    );
  }
}
