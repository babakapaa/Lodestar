import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/notifications/notification_channels.dart';
import '../../core/notifications/notifications_providers.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';

class ChannelsSection extends ConsumerWidget {
  const ChannelsSection({super.key});

  Future<void> _save(
    BuildContext context,
    WidgetRef ref,
    NotificationChannels channels,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await ref.read(notificationsRepositoryProvider).setChannels(channels);
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(notificationChannelsProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final language = Localizations.localeOf(context).languageCode;
    final channels = ref.watch(notificationChannelsProvider);
    final types = ref.watch(notificationTypesProvider).value ?? const {};
    return Section(
      title: l10n.notificationChannelsTitle,
      child: switch (channels) {
        AsyncData(:final value) when value.isEmpty => Text(
          l10n.notificationChannelsEmpty,
        ),
        AsyncData(:final value) => Column(
          children: [
            for (final item in value)
              Padding(
                padding: const EdgeInsets.only(bottom: Spacing.m),
                child: LodestarCard(
                  padding: const EdgeInsets.symmetric(
                    horizontal: Spacing.l,
                    vertical: Spacing.s,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.only(top: Spacing.s),
                        child: Text(
                          types[item.type]?.name(language) ?? item.type,
                          style: Theme.of(context).textTheme.titleSmall,
                        ),
                      ),
                      SwitchListTile(
                        contentPadding: EdgeInsets.zero,
                        title: Text(l10n.channelInApp),
                        subtitle: item.required
                            ? Text(l10n.channelRequired)
                            : null,
                        value: item.inApp,
                        onChanged: item.required
                            ? null
                            : (on) =>
                                  _save(context, ref, item.copyWith(inApp: on)),
                      ),
                      SwitchListTile(
                        contentPadding: EdgeInsets.zero,
                        title: Text(l10n.channelPush),
                        value: item.push,
                        onChanged: (on) =>
                            _save(context, ref, item.copyWith(push: on)),
                      ),
                      SwitchListTile(
                        contentPadding: EdgeInsets.zero,
                        title: Text(l10n.channelEmail),
                        value: item.email,
                        onChanged: (on) =>
                            _save(context, ref, item.copyWith(email: on)),
                      ),
                    ],
                  ),
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
