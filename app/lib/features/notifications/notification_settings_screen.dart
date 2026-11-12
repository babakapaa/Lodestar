import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'channels_section.dart';
import 'quiet_hours_section.dart';

class NotificationSettingsScreen extends StatelessWidget {
  const NotificationSettingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.notificationSettingsTitle)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: ListView(
              padding: const EdgeInsets.all(Spacing.l),
              children: const [
                QuietHoursSection(),
                SizedBox(height: Spacing.xl),
                ChannelsSection(),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
