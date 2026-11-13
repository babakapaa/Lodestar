import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'award.dart';
import 'award_dialog.dart';
import 'league_providers.dart';

class AwardsScreen extends ConsumerWidget {
  const AwardsScreen({super.key});

  Future<void> _edit(
    BuildContext context,
    WidgetRef ref, [
    Award? award,
  ]) async {
    final input = await AwardDialog.show(context, award);
    if (input == null || !context.mounted) {
      return;
    }
    final api = ref.read(leagueApiProvider);
    await TeamAction.run(
      context,
      () => award == null
          ? api.createAward(input)
          : api.updateAward(award.id, input),
    );
    ref.invalidate(awardsProvider);
  }

  List<Widget> _body(BuildContext context, WidgetRef ref, List<Award> awards) {
    final l10n = context.l10n;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    return [
      Text(l10n.awardsHint, style: muted),
      const SizedBox(height: Spacing.l),
      if (awards.isEmpty) Text(l10n.leagueAwardsNone),
      for (final award in awards)
        LodestarCard(
          key: ValueKey('award-${award.id}'),
          child: ListTile(
            contentPadding: EdgeInsets.zero,
            leading: Icon(
              award.archived
                  ? Icons.inventory_2_outlined
                  : Icons.emoji_events_outlined,
            ),
            title: Text(award.title),
            subtitle: Text(
              [
                if (award.gpuHours > 0) l10n.leagueGpuHours(award.gpuHours),
                award.quantity == null
                    ? l10n.awardGranted(award.granted)
                    : l10n.awardGrantedOf(award.granted, award.quantity!),
                if (award.archived) l10n.awardArchived,
              ].join(', '),
            ),
            onTap: () => _edit(context, ref, award),
          ),
        ),
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.awardsTitle)),
      floatingActionButton: FloatingActionButton.extended(
        key: const ValueKey('award-new'),
        onPressed: () => _edit(context, ref),
        icon: const Icon(Icons.add),
        label: Text(l10n.awardNew),
      ),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: ref.watch(awardsProvider),
          onRefresh: () => ref.refresh(awardsProvider.future),
          builder: (awards) => _body(context, ref, awards),
        ),
      ),
    );
  }
}
