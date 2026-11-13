import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/async_list_view.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'export_menu.dart';
import 'hours_dialog.dart';
import 'quota_card.dart';
import 'run_tile.dart';
import 'sandbox_providers.dart';
import 'sandbox_quota.dart';
import 'sandbox_refresh.dart';
import 'team_pool.dart';
import 'team_pool_tile.dart';

class SandboxScreen extends ConsumerWidget {
  const SandboxScreen({super.key});

  Future<void> _give(BuildContext context, WidgetRef ref, TeamPool pool) async {
    final l10n = context.l10n;
    final request = await HoursDialog.show(
      context,
      title: l10n.sandboxGiveTitle(pool.name),
      action: l10n.sandboxGive,
    );
    if (request == null || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(sandboxApiProvider).transfer(pool.id, request.hours),
    );
    ref.invalidate(sandboxQuotaProvider);
  }

  List<Widget> _body(BuildContext context, WidgetRef ref, SandboxQuota quota) {
    final l10n = context.l10n;
    if (quota.term == null) {
      return [
        EmptyState(
          icon: Icons.memory,
          title: l10n.sandboxNoTerm,
          body: l10n.sandboxNoTermBody,
        ),
      ];
    }
    final runs = ref.watch(sandboxRunsProvider).value ?? const [];
    final active = [
      for (final r in runs)
        if (r.status.active) r,
    ];
    final past = [
      for (final r in runs)
        if (!r.status.active) r,
    ];
    return [
      QuotaCard(quota: quota),
      if (quota.teams.isNotEmpty) ...[
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.sandboxTeamPools,
          child: Column(
            children: [
              for (final pool in quota.teams)
                TeamPoolTile(
                  key: ValueKey('pool-${pool.id}'),
                  pool: pool,
                  onGive: () => _give(context, ref, pool),
                ),
            ],
          ),
        ),
      ],
      const SizedBox(height: Spacing.xl),
      if (runs.isEmpty)
        EmptyState(
          icon: Icons.rocket_launch_outlined,
          title: l10n.sandboxNoRuns,
          body: l10n.sandboxNoRunsBody,
        ),
      if (active.isNotEmpty)
        Section(
          title: l10n.sandboxActive,
          child: Column(
            children: [
              for (final run in active)
                RunTile(key: ValueKey(run.id), run: run),
            ],
          ),
        ),
      if (past.isNotEmpty) ...[
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.sandboxHistory,
          child: Column(
            children: [
              for (final run in past) RunTile(key: ValueKey(run.id), run: run),
            ],
          ),
        ),
      ],
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final quota = ref.watch(sandboxQuotaProvider);
    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.sandboxTitle),
        actions: [
          IconButton(
            key: const ValueKey('sandbox-topups'),
            tooltip: l10n.sandboxTopupsTitle,
            icon: const Icon(Icons.volunteer_activism_outlined),
            onPressed: () => context.push(AppRoutes.sandboxTopups),
          ),
          const ExportMenu(),
        ],
      ),
      floatingActionButton: quota.value?.term == null
          ? null
          : FloatingActionButton.extended(
              key: const ValueKey('sandbox-new'),
              onPressed: () => context.push(AppRoutes.sandboxNew),
              icon: const Icon(Icons.play_arrow),
              label: Text(l10n.sandboxStart),
            ),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: quota,
          onRefresh: () async {
            SandboxRefresh.all(ref);
            await ref.read(sandboxQuotaProvider.future);
          },
          builder: (q) => _body(context, ref, q),
        ),
      ),
    );
  }
}
