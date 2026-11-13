import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../l10n/l10n_context.dart';
import 'league_providers.dart';
import 'league_refresh.dart';
import 'league_task_form.dart';

class LeagueTaskEditScreen extends ConsumerWidget {
  const LeagueTaskEditScreen({
    this.kind,
    this.originId,
    this.taskId,
    super.key,
  });

  final String? kind;
  final String? originId;
  final String? taskId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final id = taskId;
    return Scaffold(
      appBar: AppBar(
        title: Text(id == null ? l10n.leagueTaskNew : l10n.leagueTaskEdit),
      ),
      body: SafeArea(
        top: false,
        child: id == null
            ? AsyncListView(
                value: ref.watch(leagueProvider((kind!, originId!))),
                onRefresh: () =>
                    ref.refresh(leagueProvider((kind!, originId!)).future),
                builder: (view) => [
                  LeagueTaskForm(
                    leagueId: view.id,
                    onSaved: () => ref.refreshLeague(view.kind, view.originId),
                  ),
                ],
              )
            : AsyncListView(
                value: ref.watch(leagueTaskProvider(id)),
                onRefresh: () => ref.refresh(leagueTaskProvider(id).future),
                builder: (details) => [
                  LeagueTaskForm(
                    leagueId: details.task.leagueId,
                    task: details.task,
                    onSaved: () {
                      ref
                        ..invalidate(leagueTaskProvider(id))
                        ..refreshLeague(details.kind, details.originId);
                    },
                  ),
                ],
              ),
      ),
    );
  }
}
