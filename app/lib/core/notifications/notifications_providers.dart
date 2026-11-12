import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../auth/auth_providers.dart';
import 'notification_channels.dart';
import 'notification_page.dart';
import 'notification_type_info.dart';
import 'notifications_controller.dart';
import 'notifications_repository.dart';
import 'unread_count_controller.dart';

final notificationsRepositoryProvider = Provider<NotificationsRepository>(
  (ref) => NotificationsRepository(ref.watch(apiClientProvider)),
);

final notificationTypesProvider =
    FutureProvider<Map<String, NotificationTypeInfo>>((ref) async {
      final types = await ref.watch(notificationsRepositoryProvider).types();
      return {for (final type in types) type.key: type};
    });

final notificationsProvider =
    AsyncNotifierProvider.autoDispose<
      NotificationsController,
      NotificationPage
    >(NotificationsController.new);

final unreadCountProvider = AsyncNotifierProvider<UnreadCountController, int>(
  UnreadCountController.new,
);

final notificationChannelsProvider =
    FutureProvider.autoDispose<List<NotificationChannels>>(
      (ref) => ref.watch(notificationsRepositoryProvider).channels(),
    );
