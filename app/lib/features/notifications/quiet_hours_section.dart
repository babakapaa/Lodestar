import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/settings/settings_providers.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';

class QuietHoursSection extends ConsumerWidget {
  const QuietHoursSection({super.key});

  static const enabledKey = 'core.quiet_hours';
  static const fromKey = 'core.quiet_from';
  static const toKey = 'core.quiet_to';

  Future<void> _set(
    BuildContext context,
    WidgetRef ref,
    String key,
    Object value,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await ref.read(preferencesRepositoryProvider).set(key, value);
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(serverPreferencesProvider);
  }

  Future<void> _pickTime(
    BuildContext context,
    WidgetRef ref,
    String key,
    int minutes,
  ) async {
    final picked = await showTimePicker(
      context: context,
      initialTime: TimeOfDay(hour: minutes ~/ 60, minute: minutes % 60),
    );
    if (picked != null && context.mounted) {
      await _set(context, ref, key, picked.hour * 60 + picked.minute);
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final values = ref.watch(serverPreferencesProvider);
    return Section(
      title: l10n.quietHoursTitle,
      child: switch (values) {
        AsyncData(:final value) => Builder(
          builder: (context) {
            final enabled = value[enabledKey]?.value == true;
            final from = _minutes(value[fromKey]?.value, 22 * 60);
            final to = _minutes(value[toKey]?.value, 7 * 60);
            final format = MaterialLocalizations.of(context);
            String label(int minutes) => format.formatTimeOfDay(
              TimeOfDay(hour: minutes ~/ 60, minute: minutes % 60),
              alwaysUse24HourFormat: MediaQuery.alwaysUse24HourFormatOf(
                context,
              ),
            );
            return Column(
              children: [
                SwitchListTile(
                  contentPadding: EdgeInsets.zero,
                  title: Text(l10n.quietHoursEnabled),
                  subtitle: Text(l10n.quietHoursBody),
                  value: enabled,
                  onChanged: (on) => _set(context, ref, enabledKey, on),
                ),
                ListTile(
                  contentPadding: EdgeInsets.zero,
                  enabled: enabled,
                  title: Text(l10n.quietHoursFrom),
                  trailing: Text(label(from)),
                  onTap: () => _pickTime(context, ref, fromKey, from),
                ),
                ListTile(
                  contentPadding: EdgeInsets.zero,
                  enabled: enabled,
                  title: Text(l10n.quietHoursTo),
                  trailing: Text(label(to)),
                  onTap: () => _pickTime(context, ref, toKey, to),
                ),
              ],
            );
          },
        ),
        AsyncError(:final error) => FormError(
          message: ApiErrorText.of(l10n, error),
        ),
        _ => const Center(child: CircularProgressIndicator()),
      },
    );
  }

  static int _minutes(Object? value, int fallback) =>
      value is num ? value.toInt() : fallback;
}
