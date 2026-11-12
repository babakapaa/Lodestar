import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/notifications/app_notification.dart';
import '../../core/notifications/notification_type_info.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class NotificationTile extends StatelessWidget {
  const NotificationTile({
    required this.notification,
    required this.type,
    required this.onTap,
    super.key,
  });

  final AppNotification notification;
  final NotificationTypeInfo? type;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final palette = context.tokens.palette;
    final locale = Localizations.localeOf(context);
    final text = type?.render(notification.params, locale.languageCode);
    final time = DateFormat.MMMd(locale.toLanguageTag())
        .add_Hm()
        .format(notification.createdAt.toLocal());
    final textTheme = Theme.of(context).textTheme;
    return ListTile(
      onTap: onTap,
      leading: Icon(
        notification.read
            ? Icons.notifications_none
            : Icons.notifications_active,
        color: notification.read ? palette.muted : palette.accent,
      ),
      title: Text(
        text?.title ?? context.l10n.notificationFallback,
        style: notification.read
            ? textTheme.bodyLarge
            : textTheme.bodyLarge!.copyWith(fontWeight: FontWeight.w600),
      ),
      subtitle: Text(
        [if (text != null && text.body.isNotEmpty) text.body, time].join('\n'),
      ),
      isThreeLine: text != null && text.body.isNotEmpty,
    );
  }
}
