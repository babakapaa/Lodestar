import 'app_notification.dart';

class NotificationPage {
  const NotificationPage({required this.items, required this.hasMore});

  static const size = 50;

  final List<AppNotification> items;
  final bool hasMore;
}
