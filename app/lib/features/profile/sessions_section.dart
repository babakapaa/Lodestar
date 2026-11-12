import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'device_session.dart';
import 'profile_providers.dart';

class SessionsSection extends ConsumerWidget {
  const SessionsSection({super.key});

  Future<void> _revoke(
    BuildContext context,
    WidgetRef ref,
    DeviceSession session,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await ref.read(profileRepositoryProvider).revoke(session.id);
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(sessionsProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final sessions = ref.watch(sessionsProvider);
    final format = DateFormat.yMMMd(
      Localizations.localeOf(context).toLanguageTag(),
    ).add_Hm();
    return Section(
      title: l10n.sessionsTitle,
      child: switch (sessions) {
        AsyncData(:final value) => Column(
          children: [
            for (final session in value)
              ListTile(
                contentPadding: EdgeInsets.zero,
                leading: const Icon(Icons.devices_outlined),
                title: Text(session.deviceName),
                subtitle: Text(
                  session.current
                      ? l10n.sessionsCurrent
                      : l10n.sessionsLastUsed(
                          format.format(session.lastUsedAt.toLocal()),
                        ),
                ),
                trailing: session.current
                    ? null
                    : IconButton(
                        tooltip: l10n.sessionsRevoke,
                        icon: const Icon(Icons.logout),
                        onPressed: () => _revoke(context, ref, session),
                      ),
              ),
          ],
        ),
        AsyncError(:final error) => FormError(
          message: ApiErrorText.of(l10n, error),
        ),
        _ => const Center(child: CircularProgressIndicator()),
      },
    );
  }
}
