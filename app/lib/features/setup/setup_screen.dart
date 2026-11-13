import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/async_list_view.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'setup_overview.dart';
import 'setup_problems.dart';
import 'setup_providers.dart';
import 'setup_section_tile.dart';

class SetupScreen extends ConsumerWidget {
  const SetupScreen({super.key});

  Future<void> _launch(BuildContext context, WidgetRef ref) async {
    await TeamAction.run(context, () => ref.read(setupApiProvider).launch());
    ref
      ..invalidate(setupOverviewProvider)
      ..invalidate(setupStatusProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.setupTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView<SetupOverview>(
          value: ref.watch(setupOverviewProvider),
          onRefresh: () => ref.refresh(setupOverviewProvider.future),
          builder: (overview) {
            final done = overview.sections.where((s) => s.isApproved).length;
            final total = overview.sections.length;
            final problems = overview.problems;
            return [
              if (overview.launched) ...[
                Text(
                  l10n.setupLaunched,
                  style: Theme.of(context).textTheme.titleMedium,
                ),
                const SizedBox(height: Spacing.l),
              ],
              Text(l10n.setupProgress(done, total)),
              const SizedBox(height: Spacing.s),
              LinearProgressIndicator(value: total == 0 ? 0 : done / total),
              const SizedBox(height: Spacing.l),
              for (final section in overview.sections)
                SetupSectionTile(
                  key: ValueKey('setup-${section.key}'),
                  section: section,
                  onTap: () =>
                      context.push(AppRoutes.setupSection(section.key)),
                ),
              if (problems != null && !overview.launched) ...[
                const SizedBox(height: Spacing.xl),
                SetupProblems(problems: problems),
                const SizedBox(height: Spacing.l),
                Align(
                  alignment: Alignment.centerLeft,
                  child: LodestarButton(
                    key: const ValueKey('setup-launch-action'),
                    label: l10n.setupLaunchAction,
                    icon: Icons.rocket_launch_outlined,
                    onPressed: problems.isEmpty
                        ? () => _launch(context, ref)
                        : null,
                  ),
                ),
              ],
            ];
          },
        ),
      ),
    );
  }
}
