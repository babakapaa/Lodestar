import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../live/live_event.dart';
import '../live/live_providers.dart';
import 'app_notification.dart';
import 'notification_page.dart';
import 'notifications_providers.dart';

class NotificationsController extends AsyncNotifier<NotificationPage> {
  @override
  Future<NotificationPage> build() async {
    ref.listen(liveEventsProvider, (previous, next) {
      if (next case AsyncData(:final value)) {
        _receive(value);
      }
    });
    return _page(await ref.read(notificationsRepositoryProvider).page());
  }

  Future<void> loadMore() async {
    final current = state.value;
    if (current == null || !current.hasMore || current.items.isEmpty) {
      return;
    }
    final next = await ref
        .read(notificationsRepositoryProvider)
        .page(before: current.items.last.createdAt);
    final known = {for (final item in current.items) item.id};
    state = AsyncData(
      NotificationPage(
        items: [
          ...current.items,
          ...next.where((item) => !known.contains(item.id)),
        ],
        hasMore: next.length >= NotificationPage.size,
      ),
    );
  }

  Future<void> markRead(AppNotification notification) async {
    if (notification.read) {
      return;
    }
    await ref.read(notificationsRepositoryProvider).markRead(notification.id);
    _update((item) => item.id == notification.id);
  }

  Future<void> markAllRead() async {
    await ref.read(notificationsRepositoryProvider).markAllRead();
    _update((item) => true);
  }

  void _update(bool Function(AppNotification item) matches) {
    final current = state.value;
    if (current != null) {
      final now = DateTime.now().toUtc();
      state = AsyncData(
        NotificationPage(
          items: [
            for (final item in current.items)
              matches(item) ? item.markedRead(now) : item,
          ],
          hasMore: current.hasMore,
        ),
      );
    }
    ref.invalidate(unreadCountProvider);
  }

  void _receive(LiveEvent event) {
    if (event.type == LiveEvent.resync) {
      ref.invalidateSelf();
      return;
    }
    final payload = event.payload;
    final current = state.value;
    if (event.type != LiveEvent.notification ||
        payload is! Map<String, Object?> ||
        current == null) {
      return;
    }
    final notification = AppNotification.fromJson(payload);
    if (current.items.any((item) => item.id == notification.id)) {
      return;
    }
    state = AsyncData(
      NotificationPage(
        items: [notification, ...current.items],
        hasMore: current.hasMore,
      ),
    );
  }

  static NotificationPage _page(List<AppNotification> items) =>
      NotificationPage(
        items: items,
        hasMore: items.length >= NotificationPage.size,
      );
}
