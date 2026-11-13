import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../l10n/l10n_context.dart';
import 'sandbox_labels_context.dart';
import 'topup.dart';

class TopupTile extends StatelessWidget {
  const TopupTile({
    required this.topup,
    this.onApprove,
    this.onDecline,
    super.key,
  });

  final Topup topup;
  final VoidCallback? onApprove;
  final VoidCallback? onDecline;

  String _status(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    return switch (topup.status) {
      'approved' => l10n.sandboxTopupApproved(
        labels.hours(topup.granted ?? topup.hours),
      ),
      'declined' => l10n.sandboxTopupDeclined,
      _ => l10n.sandboxTopupPending,
    };
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    final decide = onApprove != null && onDecline != null;
    return LodestarCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ListTile(
            contentPadding: EdgeInsets.zero,
            title: Text(
              topup.mine ? topup.title : '${topup.displayName}, ${topup.title}',
            ),
            subtitle: Text(
              [
                labels.hours(topup.hours),
                if (topup.reason.isNotEmpty) topup.reason,
                if (!decide) _status(context),
              ].join('\n'),
            ),
          ),
          if (decide)
            OverflowBar(
              alignment: MainAxisAlignment.end,
              children: [
                TextButton(
                  key: ValueKey('decline-${topup.id}'),
                  onPressed: onDecline,
                  child: Text(l10n.sandboxDecline),
                ),
                FilledButton(
                  key: ValueKey('approve-${topup.id}'),
                  onPressed: onApprove,
                  child: Text(l10n.sandboxApprove),
                ),
              ],
            ),
        ],
      ),
    );
  }
}
