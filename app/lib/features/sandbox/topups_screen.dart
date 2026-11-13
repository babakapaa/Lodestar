import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'hours_dialog.dart';
import 'sandbox_providers.dart';
import 'topup.dart';
import 'topup_tile.dart';

class TopupsScreen extends ConsumerWidget {
  const TopupsScreen({super.key});

  Future<void> _approve(BuildContext context, WidgetRef ref, Topup t) async {
    final l10n = context.l10n;
    final request = await HoursDialog.show(
      context,
      title: l10n.sandboxApproveTitle(t.displayName),
      action: l10n.sandboxApprove,
      hours: t.hours,
    );
    if (request == null || !context.mounted) {
      return;
    }
    await _decide(context, ref, t, approve: true, hours: request.hours);
  }

  Future<void> _decide(
    BuildContext context,
    WidgetRef ref,
    Topup t, {
    required bool approve,
    double? hours,
  }) async {
    await TeamAction.run(
      context,
      () => ref
          .read(sandboxApiProvider)
          .decideTopup(t.id, approve: approve, hours: hours),
    );
    ref.invalidate(sandboxTopupsProvider);
  }

  List<Widget> _body(BuildContext context, WidgetRef ref, List<Topup> all) {
    final l10n = context.l10n;
    final incoming = [
      for (final t in all)
        if (!t.mine && t.pending) t,
    ];
    final mine = [
      for (final t in all)
        if (t.mine) t,
    ];
    if (incoming.isEmpty && mine.isEmpty) {
      return [
        EmptyState(
          icon: Icons.volunteer_activism_outlined,
          title: l10n.sandboxTopupsEmpty,
          body: l10n.sandboxTopupsEmptyBody,
        ),
      ];
    }
    return [
      if (incoming.isNotEmpty) ...[
        Section(
          title: l10n.sandboxTopupsIncoming,
          child: Column(
            children: [
              for (final t in incoming)
                TopupTile(
                  key: ValueKey(t.id),
                  topup: t,
                  onApprove: () => _approve(context, ref, t),
                  onDecline: () => _decide(context, ref, t, approve: false),
                ),
            ],
          ),
        ),
        const SizedBox(height: Spacing.xl),
      ],
      if (mine.isNotEmpty)
        Section(
          title: l10n.sandboxTopupsMine,
          child: Column(
            children: [
              for (final t in mine) TopupTile(key: ValueKey(t.id), topup: t),
            ],
          ),
        ),
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.sandboxTopupsTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: ref.watch(sandboxTopupsProvider),
          onRefresh: () => ref.refresh(sandboxTopupsProvider.future),
          builder: (all) => _body(context, ref, all),
        ),
      ),
    );
  }
}
