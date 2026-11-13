import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigator/journal_api.dart';
import '../../core/navigator/journal_entry.dart';
import '../../core/navigator/journal_providers.dart';
import '../../core/settings/settings_providers.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';

class JournalScreen extends ConsumerWidget {
  const JournalScreen({super.key});

  Future<void> _toggle(BuildContext context, WidgetRef ref, bool on) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await ref
          .read(preferencesRepositoryProvider)
          .set(JournalApi.preference, on);
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(serverPreferencesProvider);
  }

  Future<void> _delete(
    BuildContext context,
    WidgetRef ref,
    JournalEntry entry,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await ref.read(journalApiProvider).delete(entry.id);
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(journalProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final locale = Localizations.localeOf(context).toLanguageTag();
    final on =
        ref
            .watch(serverPreferencesProvider)
            .value?[JournalApi.preference]
            ?.value !=
        false;
    final entries = ref.watch(journalProvider);
    return Scaffold(
      appBar: AppBar(title: Text(l10n.journalTitle)),
      body: SafeArea(
        top: false,
        child: RefreshIndicator(
          onRefresh: () => ref.refresh(journalProvider.future),
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              SwitchListTile(
                contentPadding: EdgeInsets.zero,
                title: Text(l10n.journalEnabled),
                subtitle: Text(l10n.journalHint),
                value: on,
                onChanged: (value) => _toggle(context, ref, value),
              ),
              const SizedBox(height: Spacing.l),
              switch (entries) {
                AsyncData(:final value) when value.isEmpty => EmptyState(
                  icon: Icons.history,
                  title: l10n.journalEmpty,
                  body: l10n.journalHint,
                ),
                AsyncData(:final value) => Column(
                  children: [
                    for (final entry in value.reversed)
                      ListTile(
                        key: ValueKey(entry.id),
                        contentPadding: EdgeInsets.zero,
                        title: Text(entry.query),
                        subtitle: Text(
                          [
                            DateFormat.yMMMd(locale)
                                .add_Hm()
                                .format(entry.createdAt.toLocal()),
                            ...entry.sources,
                          ].join('\n'),
                        ),
                        trailing: IconButton(
                          tooltip: l10n.journalDelete,
                          icon: const Icon(Icons.delete_outline),
                          onPressed: () => _delete(context, ref, entry),
                        ),
                      ),
                  ],
                ),
                AsyncError(:final error) => FormError(
                  message: ApiErrorText.of(l10n, error),
                ),
                _ => const Center(child: CircularProgressIndicator()),
              },
            ],
          ),
        ),
      ),
    );
  }
}
