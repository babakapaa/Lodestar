import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/navigation/app_routes.dart';
import '../../core/notifications/app_notification.dart';
import '../../core/notifications/notifications_providers.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'notification_tile.dart';

class NotificationsScreen extends ConsumerWidget {
  const NotificationsScreen({super.key});

  Future<void> _run(
    BuildContext context,
    Future<void> Function() action,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await action();
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final page = ref.watch(notificationsProvider);
    final types = ref.watch(notificationTypesProvider).value ?? const {};
    final controller = ref.read(notificationsProvider.notifier);
    final hasUnread =
        page.value?.items.any((AppNotification item) => !item.read) ?? false;
    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.notificationsTitle),
        actions: [
          if (hasUnread)
            IconButton(
              tooltip: l10n.notificationsReadAll,
              onPressed: () => _run(context, controller.markAllRead),
              icon: const Icon(Icons.done_all),
            ),
          IconButton(
            tooltip: l10n.notificationSettingsTitle,
            onPressed: () => context.push(AppRoutes.notificationSettings),
            icon: const Icon(Icons.tune),
          ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: switch (page) {
              AsyncData(:final value) when value.items.isEmpty => EmptyState(
                icon: Icons.notifications_none,
                title: l10n.notificationsEmptyTitle,
                body: l10n.notificationsEmptyBody,
              ),
              AsyncData(:final value) => RefreshIndicator(
                onRefresh: () => ref.refresh(notificationsProvider.future),
                child: ListView(
                  children: [
                    for (final item in value.items)
                      NotificationTile(
                        notification: item,
                        type: types[item.type],
                        onTap: () =>
                            _run(context, () => controller.markRead(item)),
                      ),
                    if (value.hasMore)
                      Padding(
                        padding: const EdgeInsets.all(Spacing.l),
                        child: LodestarButton(
                          label: l10n.notificationsMore,
                          variant: LodestarButtonVariant.secondary,
                          onPressed: () => _run(context, controller.loadMore),
                        ),
                      ),
                  ],
                ),
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
