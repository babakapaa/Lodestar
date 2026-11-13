import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'navigator_intent.dart';
import 'navigator_labels_context.dart';

class IntentSheet extends StatelessWidget {
  const IntentSheet({super.key});

  static Future<NavigatorIntent?> show(BuildContext context) {
    return showModalBottomSheet<NavigatorIntent>(
      context: context,
      showDragHandle: true,
      builder: (context) => const IntentSheet(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final labels = context.navigatorLabels;
    return SafeArea(
      child: ListView(
        shrinkWrap: true,
        children: [
          ListTile(
            title: Text(
              context.l10n.navigatorPickIntent,
              style: Theme.of(context).textTheme.titleMedium,
            ),
          ),
          for (final intent in NavigatorIntent.choosable)
            ListTile(
              title: Text(labels.intent(intent)),
              onTap: () => Navigator.of(context).pop(intent),
            ),
        ],
      ),
    );
  }
}
