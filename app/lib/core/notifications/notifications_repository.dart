import '../api/api_client.dart';
import 'app_notification.dart';
import 'notification_channels.dart';
import 'notification_type_info.dart';

class NotificationsRepository {
  const NotificationsRepository(this._api);

  final ApiClient _api;

  Future<List<AppNotification>> page({DateTime? before}) async {
    final body = await _api.get('/v1/notifications', {
      if (before != null) 'before': before.toUtc().toIso8601String(),
    });
    return [
      for (final item in body! as List<Object?>)
        AppNotification.fromJson(item! as Map<String, Object?>),
    ];
  }

  Future<int> unread() async {
    final body = await _api.get('/v1/notifications/unread');
    return ((body! as Map<String, Object?>)['unread']! as num).toInt();
  }

  Future<void> markRead(String id) =>
      _api.post('/v1/notifications/${Uri.encodeComponent(id)}/read');

  Future<void> markAllRead() => _api.post('/v1/notifications/read-all');

  Future<List<NotificationTypeInfo>> types() async => [
    for (final item
        in (await _api.get('/v1/notification-types'))! as List<Object?>)
      NotificationTypeInfo.fromJson(item! as Map<String, Object?>),
  ];

  Future<List<NotificationChannels>> channels() async => [
    for (final item
        in (await _api.get('/v1/me/notification-channels'))! as List<Object?>)
      NotificationChannels.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> setChannels(NotificationChannels channels) => _api.put(
    '/v1/me/notification-channels/${Uri.encodeComponent(channels.type)}',
    channels.toInput(),
  );
}
